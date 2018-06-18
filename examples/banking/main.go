package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"

	"github.com/jmalloc/ax/examples/banking/domain"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/examples/banking/projections"
	"github.com/jmalloc/ax/examples/banking/workflows"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/projection"
	"github.com/jmalloc/ax/src/ax/routing"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/mapping/direct"
	"github.com/jmalloc/ax/src/ax/saga/mapping/keyset"
	"github.com/jmalloc/ax/src/ax/saga/persistence/crud"
	"github.com/jmalloc/ax/src/ax/saga/persistence/eventsourcing"
	"github.com/jmalloc/ax/src/axcli"
	"github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/axrmq"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

func main() {
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
			Saga: workflows.TransferWorkflow,
			Mapper: keyset.ByField(
				axmysql.SagaKeySetRepository,
				"TransferId",
			),
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

	ep := &endpoint.Endpoint{
		Name: "ax.examples.banking",
		Transport: &axrmq.Transport{
			Conn: rmq,
		},
		In: &observability.InboundHook{
			Observers: observers,
			Next: &persistence.Injector{
				DataStore: ds,
				Next: &outbox.Deduplicator{
					Repository: axmysql.OutboxRepository,
					Next: &routing.Dispatcher{
						Routes: htable,
					},
				},
			},
		},
		Out: &observability.OutboundHook{
			Observers: observers,
			Next: &routing.Router{
				Routes: etable,
				Next:   &endpoint.TransportStage{},
			},
		},
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
				return con.Consume(ctx)
			})

			return g.Wait()
		},
	})

	err = cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
