package main

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/axrmq"
	"github.com/streadway/amqp"
)

var handler = bus.MessageHandlerFunc(
	ax.TypesOf(&messages.OpenAccount{}),
	func(ctx ax.MessageContext, m ax.Message) error {
		return nil
	},
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

	dtable, err := bus.NewDispatchTable(handler)
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

	if err := xport.Initialize(ctx, "ax.examples.banking"); err != nil {
		panic(err)
	}

	if err := in.Initialize(ctx, xport); err != nil {
		panic(err)
	}

	if err := out.Initialize(ctx, xport); err != nil {
		panic(err)
	}

	for {
		m, err := xport.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		err = in.DeliverMessage(ctx, out, m)
		op := bus.OpAck
		if err == nil {
			if m.DeliveryCount < 3 {
				op = bus.OpRetry
			} else {
				op = bus.OpReject
			}
		}

		if err := m.Done(ctx, op); err != nil {
			panic(err)
		}
	}
}
