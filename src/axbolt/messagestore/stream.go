package messagestore

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax"
)

const (
	// DefaultFetchLimit is the number of messages to fetch in each select query on
	// a message stream.
	DefaultFetchLimit = 100

	// DefaultPollInterval is the default time to wait between polls in
	// MessageStream.Next().
	DefaultPollInterval = 500 * time.Millisecond
)

// Stream is a Bolt-backed implementation of Ax's messagestore.Stream
// interface.
type Stream struct {
	Fetcher      Fetcher
	NextOffset   uint64
	Limit        uint64
	PollInterval time.Duration

	envpbs map[uint64]*ax.EnvelopeProto
}

// Next advances the stream to the next message.
//
// It blocks until a message is available, or ctx is canceled.
func (s *Stream) Next(ctx context.Context) error {
	ok, err := s.TryNext(ctx)
	if ok || err != nil {
		return err
	}

	p := s.PollInterval
	if p == 0 {
		p = DefaultPollInterval
	}

	tick := time.NewTicker(p)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			ok, err := s.TryNext(ctx)
			if ok || err != nil {
				return err
			}
		}
	}
}

// TryNext advances the stream to the next message.
//
// It returns what if there are no more messages in the stream.
func (s *Stream) TryNext(ctx context.Context) (bool, error) {
	if s.envpbs != nil {
		if s.advance() {
			return true, nil
		}
	}
	if err := s.fetchMessages(ctx); err != nil {
		return false, err
	}

	return s.advance(), nil
}

// Get returns the message at the current offset in the stream.
func (s *Stream) Get(ctx context.Context) (ax.Envelope, error) {
	if s.envpbs == nil {
		panic("Next() must be called before Get()")
	}

	envproto := s.envpbs[s.NextOffset]
	return ax.NewEnvelopeFromProto(envproto)
}

// Offset returns the offset of the message returned by Get().
func (s *Stream) Offset() (uint64, error) {
	if s.envpbs == nil {
		panic("Next() must be called before Offset()")
	}
	return s.NextOffset - 1, nil
}

// Close closes the stream.
func (s *Stream) Close() error {
	s.envpbs = nil
	return nil
}

// fetchMessages selects the next batch of messages from the stream.
func (s *Stream) fetchMessages(ctx context.Context) error {
	n := s.Limit
	if n == 0 {
		n = DefaultFetchLimit
	}

	envpbs, err := s.Fetcher.FetchMessages(ctx, s.NextOffset, n)
	if err != nil {
		return err
	}
	s.envpbs = envpbs
	return nil
}

// advance moves to the next message in s.envpbs.
func (s *Stream) advance() bool {
	o := s.NextOffset
	o++
	if _, ok := s.envpbs[o]; ok {
		s.NextOffset++
		return true
	}
	return false
}
