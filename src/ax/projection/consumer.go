package projection

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Consumer reads from a stream of messages and forwards them to a projector.
type Consumer struct {
	DataStore    persistence.DataStore
	MessageStore messagestore.GloballyOrderedStore
	Offsets      OffsetStore
	Projector    Projector

	name   string
	types  ax.MessageTypeSet
	stream messagestore.Stream
}

// Run reads pipes messages from the message stream to the projector until an
// error occurs or ctx is canceled.
func (c *Consumer) Run(ctx context.Context) error {
	c.name = c.Projector.ProjectorName()
	c.types = c.Projector.MessageTypes()

	o, err := c.Offsets.LoadOffset(ctx, c.DataStore, c.name)
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
		err = c.process(ctx)
		if err != nil {
			return err
		}
	}
}

func (c *Consumer) process(ctx context.Context) error {
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
		err = c.Projector.HandleMessage(
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
		c.name,
		o,
		o+1,
	)
	if err != nil {
		return err
	}

	return com.Commit()
}
