package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/axrmq"
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

	if err := xport.Initialize(ctx, "ax.examples.banking"); err != nil {
		panic(err)
	}

	if err := in.Initialize(ctx, xport); err != nil {
		panic(err)
	}

	if err := out.Initialize(ctx, xport); err != nil {
		panic(err)
	}

	go send(out)

	for {
		env, err := xport.ReceiveMessage(ctx)
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
			panic(err)
		}
	}
}

func send(s bus.MessageSink) {
	// ctx := context.Background()

	// snd := bus.SinkSender{Sink: s}
}
