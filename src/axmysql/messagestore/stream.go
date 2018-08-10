package messagestore

import (
	"context"
	"database/sql"
	"log"
	"time"

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

// Stream is a MySQL-backed implementation of Ax's messagestore.Stream
// interface.
type Stream struct {
	Fetcher      Fetcher
	NextOffset   uint64
	Limit        uint64
	PollInterval time.Duration

	rows *sql.Rows
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
// It returns false if there are no more messages in the stream.
func (s *Stream) TryNext(ctx context.Context) (bool, error) {
	if s.rows != nil {
		if s.advance() {
			return true, nil
		}
	}

	if err := s.fetchRows(ctx); err != nil {
		return false, err
	}

	return s.advance(), nil
}

// Get returns the message at the current offset in the stream.
func (s *Stream) Get(ctx context.Context) (ax.Envelope, error) {
	if s.rows == nil {
		panic("Next() must be called before Get()")
	}

	var (
		env         ax.Envelope
		contentType string
		data        []byte
		createdAt   string
		sendAt      string
	)

	err := s.rows.Scan(
		&env.MessageID,
		&env.CausationID,
		&env.CorrelationID,
		&createdAt,
		&sendAt,
		&contentType,
		&data,
	)
	if err != nil {
		return ax.Envelope{}, err
	}

	err = marshaling.UnmarshalTime(createdAt, &env.CreatedAt)
	if err != nil {
		return ax.Envelope{}, err
	}

	err = marshaling.UnmarshalTime(sendAt, &env.SendAt)
	if err != nil {
		return ax.Envelope{}, err
	}

	env.Message, err = ax.UnmarshalMessage(contentType, data)

	log.Printf("Get(), env: %+v", env)
	return env, err
}

// Offset returns the offset of the message returned by Get().
func (s *Stream) Offset() (uint64, error) {
	if s.rows == nil {
		panic("Next() must be called before Offset()")
	}

	return s.NextOffset - 1, nil
}

// Close closes the stream.
func (s *Stream) Close() error {
	return s.replaceRows(nil)
}

// fetchRows selects the next batch of messages from the stream.
func (s *Stream) fetchRows(ctx context.Context) error {
	n := s.Limit
	if n == 0 {
		n = DefaultFetchLimit
	}

	rows, err := s.Fetcher.FetchRows(ctx, s.NextOffset, n)
	if err != nil {
		return err
	}

	return s.replaceRows(rows)
}

// replaceRows replaces s.rows with r, closing the existing s.rows value if it
// is not nil.
func (s *Stream) replaceRows(r *sql.Rows) error {
	if s.rows != nil {
		prev := s.rows
		s.rows = nil

		if err := prev.Close(); err != nil {
			return err
		}
	}

	s.rows = r
	return nil
}

// advance moves to the next row in s.rows.
func (s *Stream) advance() bool {
	if s.rows.Next() {
		s.NextOffset++
		return true
	}

	return false
}
