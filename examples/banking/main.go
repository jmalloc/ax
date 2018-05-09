package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axcli"
	"github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/axrmq"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

type handler struct{}

func (handler) MessageTypes() ax.MessageTypeSet {
	return ax.TypesOf(&messages.OpenAccount{})
}

func (handler) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope) error {
	m := env.Message.(*messages.OpenAccount)
	return s.PublishEvent(ctx, &messages.AccountOpened{
		AccountId: m.AccountId,
		Name:      m.Name,
	})
}

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

	dtable, err := bus.NewDispatchTable(handler{})
	if err != nil {
		panic(err)
	}

	rtable, err := bus.NewRoutingTable()
	if err != nil {
		panic(err)
	}

	observers := []interface{}{
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
					Next: &bus.Dispatcher{
						Routes: dtable,
					},
				},
			},
		},
		Out: &observability.OutboundHook{
			Observers: observers,
			Next: &bus.Router{
				Routes: rtable,
				Next:   &bus.TransportStage{},
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
