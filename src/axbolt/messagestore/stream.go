package messagestore

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/marshaling"
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

	msgs map[uint64]*StoredMessage
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
	if s.msgs != nil {
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
	if s.msgs == nil {
		panic("Next() must be called before Get()")
	}

	m := s.msgs[s.NextOffset-1]
	return parseMessage(m)
}

// Offset returns the offset of the message returned by Get().
func (s *Stream) Offset() (uint64, error) {
	if s.msgs == nil {
		panic("Next() must be called before Offset()")
	}

	return s.NextOffset - 1, nil
}

// Close closes the stream.
func (s *Stream) Close() error {
	s.msgs = nil
	return nil
}

// fetchMessages selects the next batch of messages from the stream.
func (s *Stream) fetchMessages(ctx context.Context) error {
	n := s.Limit
	if n == 0 {
		n = DefaultFetchLimit
	}

	msgs, err := s.Fetcher.FetchMessages(ctx, s.NextOffset, n)
	if err != nil {
		return err
	}

	s.msgs = msgs
	return nil
}

// advance moves to the next message in s.msgs.
func (s *Stream) advance() bool {
	o := s.NextOffset
	o++
	if _, ok := s.msgs[o]; ok {
		s.NextOffset++
		return true
	}
	return false
}

// parseMessage parses a message into an ax.Envelope struct
func parseMessage(m *StoredMessage) (ax.Envelope, error) {
	var (
		err error
		env ax.Envelope
	)
	var x ptypes.DynamicAny
	if err = ptypes.UnmarshalAny(m.Message, &x); err != nil {
		return env, err
	}
	env.Message, _ = x.Message.(ax.Message)

	if err = env.MessageID.Parse(m.GetId()); err != nil {
		return env, err
	}
	if err = env.CorrelationID.Parse(m.GetCorrelationId()); err != nil {
		return env, err
	}
	if err = marshaling.UnmarshalTime(m.GetCreatedAt(), &env.CreatedAt); err != nil {
		return env, err
	}
	err = marshaling.UnmarshalTime(m.GetSendAt(), &env.SendAt)
	return env, err
}
