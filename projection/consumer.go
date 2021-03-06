package projection

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/messagestore"
	"github.com/jmalloc/ax/observability"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/twelf/src/twelf"
	opentracing "github.com/opentracing/opentracing-go"
)

// GlobalStoreConsumer reads messages from all streams in a message store and
// forwards them to an application-defined projector to produce a projection.
type GlobalStoreConsumer struct {
	Projector    Projector
	DataStore    persistence.DataStore
	MessageStore messagestore.GloballyOrderedStore
	Offsets      OffsetStore
	Logger       twelf.Logger

	key    string
	types  ax.MessageTypeSet
	stream messagestore.Stream
}

// Consume reads messages from the store and forwards them to the projector until
// an error occurs or ctx is canceled.
func (c *GlobalStoreConsumer) Consume(ctx context.Context) error {
	c.key = c.Projector.PersistenceKey()
	c.types = c.Projector.MessageTypes()

	o, err := c.Offsets.LoadOffset(ctx, c.DataStore, c.key)
	if err != nil {
		return err
	}

	c.stream, err = c.MessageStore.OpenGlobal(ctx, c.DataStore, o)
	if err != nil {
		return err
	}
	defer c.stream.Close()

	ctx = persistence.WithDataStore(ctx, c.DataStore)

	for {
		err = c.processNextMessage(ctx)
		if err != nil {
			return err
		}
	}
}

func (c *GlobalStoreConsumer) processNextMessage(ctx context.Context) error {
	err := c.stream.Next(ctx)
	if err != nil {
		return err
	}

	env, err := c.stream.Get(ctx)
	if err != nil {
		return err
	}

	tx, com, err := persistence.GetOrBeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	if c.types.Has(env.Type()) {
		mctx := ax.NewMessageContext(
			env,
			opentracing.SpanFromContext(ctx),
			observability.NewProjectionLogger(
				c.Logger,
				env,
			),
		)

		err = c.Projector.ApplyMessage(
			persistence.WithTx(ctx, tx),
			mctx,
		)
		if err != nil {
			return err
		}
	}

	o, err := c.stream.Offset()
	if err != nil {
		return err
	}

	err = c.Offsets.IncrementOffset(
		ctx,
		tx,
		c.key,
		o,
	)
	if err != nil {
		return err
	}

	return com.Commit()
}
