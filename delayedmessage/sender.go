package delayedmessage

import (
	"context"
	"time"

	"github.com/jmalloc/ax/endpoint"
	"github.com/jmalloc/ax/persistence"
)

// DefaultPollInterval is the duration to wait before checking for new messages
// to send.
var DefaultPollInterval = 15 * time.Second

// state is a function that handles a single state of the sender.
type state func(ctx context.Context) (state, error)

// Sender is a service that sends delayed messages when they become ready to be
// sent.
type Sender struct {
	DataStore        persistence.DataStore
	Repository       Repository
	OutboundPipeline endpoint.OutboundPipeline
	PollInterval     time.Duration
}

// Run sends messages as they become ready to send until ctx is canceled or an
// error occurrs.
func (s *Sender) Run(ctx context.Context) error {
	for {
		if err := s.tick(ctx); err != nil {
			return err
		}
	}
}

// Tick loads the next message from the repository and sends it if it is ready
// to be sent. Otherwise it waits for the poll interval or until the message is
// ready to be sent then tries again.
func (s *Sender) tick(ctx context.Context) error {
	env, ok, err := s.Repository.LoadNextMessage(ctx, s.DataStore)
	if err != nil {
		return err
	}

	d := s.PollInterval
	if d == 0 {
		d = DefaultPollInterval
	}

	if ok {
		delay := time.Until(env.SendAt)

		if delay <= 0 {
			return s.send(ctx, env)
		} else if delay < d {
			d = delay
		}
	}

	return s.sleep(ctx, d)
}

// send sends a message and marks it as sent.
func (s *Sender) send(ctx context.Context, env endpoint.OutboundEnvelope) error {
	if err := s.OutboundPipeline.Accept(ctx, env); err != nil {
		return err
	}

	tx, com, err := s.DataStore.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	if err := s.Repository.MarkAsSent(ctx, tx, env); err != nil {
		return err
	}

	return com.Commit()
}

// sleep blocks until ctx is canceled or the given duration elapses.
func (s *Sender) sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
