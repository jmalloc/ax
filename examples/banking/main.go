package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
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

func (handler) HandleMessage(_ context.Context, _ ax.Sender, env ax.Envelope) error {
	spew.Dump(env)
	return nil
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

	in := &persistence.Injector{
		DataStore: &axmysql.DataStore{DB: db},
		Next: &bus.Dispatcher{
			Routes: dtable,
		},
	}

	out := &bus.Router{
		Routes: rtable,
		Next:   &bus.TransportStage{},
	}

	xport := &axrmq.Transport{
		Conn: rmq,
	}

	ctx := context.Background()
	ep := "ax.examples.banking"

	if err := xport.Initialize(ctx, ep); err != nil {
		panic(err)
	}

	if err := in.Initialize(ctx, xport); err != nil {
		panic(err)
	}

	if err := out.Initialize(ctx, xport); err != nil {
		panic(err)
	}

	// -------------------------------------------------------

	cli := &cobra.Command{
		Use:   "banking",
		Short: "A basic CQRS/ES example written in Ax",
	}

	commands, err := axcli.NewCommands(
		bus.SinkSender{
			Sink: out,
		},
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
		Short: fmt.Sprintf("Run the '%s' endpoint", ep),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			for {
				env, err := xport.Produce(ctx)
				if err != nil {
					panic(err)
				}

				err = in.Accept(ctx, out, env)
				op := bus.OpAck
				if err == nil {
					if env.DeliveryCount < 3 {
						op = bus.OpRetry
					} else {
						op = bus.OpReject
					}
				}

				if err := env.Done(ctx, op); err != nil {
					return err
				}
			}
		},
	})

	err = cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
