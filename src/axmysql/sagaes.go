package axmysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/eventsourcing"
)

type EventStore struct{}

func (EventStore) AppendEvents(
	ctx context.Context,
	tx persistence.Tx,
	id saga.InstanceID,
	rev saga.Revision,
	events []ax.Event,
) error {
	stx := tx.(*Tx).sqlTx

	if rev == 0 {
		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_instance SET
				instance_id = ?,
				current_revision = 1`,
			id,
		); err != nil {
			return err
		}
	} else {
		row := stx.QueryRowContext(
			ctx,
			`SELECT
				i.current_revision
			FROM saga_instance AS i
			WHERE i.instance_id = ?
			FOR UPDATE`,
			id,
		)

		var currentRev saga.Revision
		if err := row.Scan(&currentRev); err != nil {
			return err
		}

		if rev != currentRev {
			return fmt.Errorf(
				"can not update saga instance %s, revision is out of date",
				id,
			)
		}
	}

	for _, ev := range events {
		contentType, data, err := ax.MarshalMessage(ev)
		if err != nil {
			return err
		}

		rev++

		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_event SET
				instance_id = ?,
				revision = ?,
				description = ?,
				content_type = ?,
				data = ?`,
			id,
			rev,
			ev.Description(),
			contentType,
			data,
		); err != nil {
			return err
		}
	}

	_, err := stx.ExecContext(
		ctx,
		`UPDATE saga_instance SET
			current_revision = ?
		WHERE instance_id = ?`,
		rev,
		id,
	)

	return err
}

func (EventStore) OpenStream(
	ctx context.Context,
	tx persistence.Tx,
	id saga.InstanceID,
	rev saga.Revision,
) (eventsourcing.EventStream, error) {
	return &EventStream{
		tx:  tx.(*Tx).sqlTx,
		id:  id,
		rev: rev,
	}, nil
}

type EventStream struct {
	tx  *sql.Tx
	id  saga.InstanceID
	rev saga.Revision
	ev  ax.Event
}

// Next advances the stream to the next event.
// It returns false if there are no more events.
func (s *EventStream) Next(ctx context.Context) (bool, error) {
	row := s.tx.QueryRowContext(
		ctx,
		`SELECT
			e.content_type,
			e.data
		FROM saga_event AS e
		WHERE e.instance_id = ?
		AND e.revision = ?`,
		s.id,
		s.rev,
	)

	var (
		contentType string
		data        []byte
	)

	if err := row.Scan(&contentType, &data); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	ev, err := ax.UnmarshalMessage(contentType, data)
	if err != nil {
		return false, err
	}

	s.ev = ev.(ax.Event)
	s.rev++

	return true, nil
}

// Get returns the event at the current location in the stream.
func (s *EventStream) Get(ctx context.Context) (ax.Event, error) {
	if s.ev == nil {
		return nil, errors.New("nO! TODO")
	}

	return s.ev, nil
}

// Close closes the stream.
func (s *EventStream) Close() error {
	return nil
}
