package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axcli"
	"github.com/jmalloc/ax/axmysql"
	"github.com/jmalloc/ax/axrmq"
	"github.com/jmalloc/ax/delayedmessage"
	"github.com/jmalloc/ax/endpoint"
	"github.com/jmalloc/ax/examples/banking/domain"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/examples/banking/projections"
	"github.com/jmalloc/ax/examples/banking/workflows"
	"github.com/jmalloc/ax/observability"
	"github.com/jmalloc/ax/outbox"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/projection"
	"github.com/jmalloc/ax/routing"
	"github.com/jmalloc/ax/saga"
	"github.com/jmalloc/ax/saga/mapping/direct"
	"github.com/jmalloc/ax/saga/mapping/keyset"
	"github.com/jmalloc/ax/saga/persistence/crud"
	"github.com/jmalloc/ax/saga/persistence/eventsourcing"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"github.com/uber/jaeger-client-go/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	os.Exit(run()) // os.Exit bypasses defer statements, perform them in run() instead
}

func run() int {
	db, err := sql.Open("mysql", os.Getenv("AX_MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rmq, err := amqp.Dial(os.Getenv("AX_RMQ_DSN"))
	if err != nil {
		panic(err)
	}
	defer rmq.Close()

	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	crudPersister := &crud.Persister{
		Repository: axmysql.SagaCRUDRepository,
	}

	esPersister := &eventsourcing.Persister{
		MessageStore:      axmysql.MessageStore,
		Snapshots:         axmysql.SagaSnapshotRepository,
		SnapshotFrequency: 3,
	}

	htable, err := routing.NewHandlerTable(
		// aggregates ...
		&saga.MessageHandler{
			Saga:      saga.NewAggregate(&domain.Account{}),
			Mapper:    direct.ByField("AccountId"),
			Persister: esPersister,
		},
		&saga.MessageHandler{
			Saga:      saga.NewAggregate(&domain.Transfer{}),
			Mapper:    direct.ByField("TransferId"),
			Persister: esPersister,
		},

		// workflows ...
		&saga.MessageHandler{
			Saga:      saga.NewWorkflow(&workflows.Transfer{}),
			Mapper:    keyset.ByField(axmysql.SagaKeySetRepository, "TransferId"),
			Persister: crudPersister,
		},
	)
	if err != nil {
		panic(err)
	}

	etable, err := routing.NewEndpointTable()
	if err != nil {
		panic(err)
	}

	observers := []observability.Observer{
		&observability.LoggingObserver{},
	}

	ds := axmysql.NewDataStore(db)

	// the router is the point within the outbound pipeline that is shared between
	// the delayed message sender and the endpoint itself.
	router := &routing.Router{
		Routes: etable,
		Next:   &endpoint.TransportStage{},
	}

	dms := &delayedmessage.Sender{
		DataStore:  ds,
		Repository: axmysql.DelayedMessageRepository,
		OutboundPipeline: endpoint.OutboundTracer{
			Tracer: tracer,
			Next: &observability.OutboundHook{
				Observers: observers,
				Next:      router,
			},
		},
	}

	transport := &axrmq.Transport{
		Conn:   rmq,
		Tracer: tracer,
	}

	ep := &endpoint.Endpoint{
		Name:              "ax.examples.banking",
		InboundTransport:  transport,
		OutboundTransport: transport,
		InboundPipeline: &observability.InboundHook{
			Observers: observers,
			Next: &persistence.InboundInjector{
				DataStore: ds,
				Next: &outbox.Deduplicator{
					Repository: axmysql.OutboxRepository,
					Next: &routing.Dispatcher{
						Routes: htable,
					},
				},
			},
		},
		OutboundPipeline: endpoint.OutboundTracer{
			Tracer: tracer,
			Next: &persistence.OutboundInjector{
				DataStore: ds,
				Next: &observability.OutboundHook{
					Observers: observers,
					Next: &delayedmessage.Interceptor{
						Repository: axmysql.DelayedMessageRepository,
						Next:       router,
					},
				},
			},
		},
		Tracer: tracer,
	}

	con := &projection.GlobalStoreConsumer{
		Projector:    projections.AccountProjector,
		DataStore:    ds,
		MessageStore: axmysql.MessageStore,
		Offsets:      axmysql.ProjectionOffsetStore,
	}

	// -------------------------------------------------------

	cli := &cobra.Command{
		Use:   "banking",
		Short: "A basic CQRS/ES example written in Ax",
	}

	ctx := context.Background()

	sender, err := ep.NewSender(ctx)
	if err != nil {
		panic(err)
	}

	commands, err := axcli.NewCommands(
		sender,
		ax.TypesOf(
			&messages.OpenAccount{},
			&messages.CreditAccount{},
			&messages.DebitAccount{},

			&messages.StartTransfer{},
		),
	)
	if err != nil {
		panic(err)
	}

	cli.AddCommand(commands...)
	cli.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: fmt.Sprintf("Run the '%s' endpoint", ep.Name),
		RunE: func(*cobra.Command, []string) error {
			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return ep.StartReceiving(ctx)
			})

			g.Go(func() error {
				return dms.Run(ctx)
			})

			g.Go(func() error {
				return con.Consume(ctx)
			})

			return g.Wait()
		},
	})

	err = cli.Execute()
	if err != nil {
		return 1
	}

	return 0
}
