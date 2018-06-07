package projection

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// GlobalStoreConsumer reads messages from all streams in a message store and
// forwards them to an application-defined projector to produce a projection.
type GlobalStoreConsumer struct {
	Projector    Projector
	DataStore    persistence.DataStore
	MessageStore messagestore.GloballyOrderedStore
	Offsets      OffsetStore

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
		err = c.Projector.ApplyMessage(
			persistence.WithTx(ctx, tx),
			env,
		)
		if err != nil {
			return err
		}
	}

	o, err := c.stream.Offset()
	if err != nil {
		return err
	}

	err = c.Offsets.SaveOffset(
		ctx,
		tx,
		c.key,
		o,
		o+1,
	)
	if err != nil {
		return err
	}

	return com.Commit()
}
