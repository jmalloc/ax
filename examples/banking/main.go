package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmalloc/ax/examples/banking/account"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/routing"
	"github.com/jmalloc/ax/src/ax/saga/eventsourcing"
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

	htable, err := routing.NewHandlerTable(
		&eventsourcing.MessageHandler{
			Saga:   account.AggregateRoot,
			Mapper: &axmysql.SagaMapper{},
			Repository: &eventsourcing.MessageStoreRepository{
				MessageStore: &axmysql.MessageStore{},
			},
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

	ep := &endpoint.Endpoint{
		Name: "ax.examples.banking",
		Transport: &axrmq.Transport{
			Conn: rmq,
		},
		In: &observability.InboundHook{
			Observers: observers,
			Next: &persistence.Injector{
				DataStore: &axmysql.DataStore{DB: db},
				Next: &outbox.Deduplicator{
					Repository: &axmysql.OutboxRepository{},
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
			return ep.StartReceiving(ctx)
		},
	})

	err = cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
