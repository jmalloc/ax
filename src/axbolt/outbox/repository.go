package outbox

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"

	"github.com/coreos/bbolt"

	"github.com/jmalloc/ax/src/ax/endpoint"

	"github.com/golang/protobuf/proto"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// Repository is a Bolt-backed implementation of Ax's outbox.Repository
// interface.
type Repository struct{}

const (
	// outboxBktName is the name of of the Bolt root bucket where all
	// outbox-specific data is stored.
	outboxBktName = "ax_outbox"
)

// LoadOutbox loads the unsent outbound messages that were produced when the
// message identified by id was first delivered.
func (Repository) LoadOutbox(
	ctx context.Context,
	ds persistence.DataStore,
	id ax.MessageID,
) ([]endpoint.OutboundEnvelope, bool, error) {
	db := boltpersistence.ExtractDB(ds)
	tx, err := db.Begin(false)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	bkt := boltutil.GetBkt(
		tx,
		outboxBktName,
		id.Get(),
	)
	if bkt == nil {
		return nil, false, nil
	}

	c := bkt.Cursor()

	var envelopes []endpoint.OutboundEnvelope
	for k, v := c.First(); k != nil; k, v = c.Next() {
		env, err := parseOutboxMessage(v, id)
		if err != nil {
			return nil, false, err
		}
		envelopes = append(envelopes, env)
	}

	return envelopes, true, nil
}

// ErrOutboxExists is returned by SaveOutbox method if the outbox messages
// with the same causation id already exist in the database
var ErrOutboxExists = errors.New("outbox already exists in the database")

// SaveOutbox saves a set of unsent outbound messages that were produced
// when the message identified by id was delivered.
func (Repository) SaveOutbox(
	ctx context.Context,
	ptx persistence.Tx,
	id ax.MessageID,
	envs []endpoint.OutboundEnvelope,
) error {
	tx := boltpersistence.ExtractTx(ptx)

	if b := boltutil.GetBkt(
		tx,
		outboxBktName,
		id.Get(),
	); b != nil {
		return ErrOutboxExists
	}

	bkt, err := boltutil.MakeBkt(
		tx,
		outboxBktName,
		id.Get(),
	)
	if err != nil {
		return err
	}

	for _, env := range envs {
		if err := insertOutboxMessage(bkt, env); err != nil {
			return err
		}
	}

	return nil
}

// MarkAsSent marks a message as sent, removing it from the outbox.
func (Repository) MarkAsSent(
	ctx context.Context,
	ptx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {
	tx := boltpersistence.ExtractTx(ptx)

	bkt := boltutil.GetBkt(
		tx,
		outboxBktName,
		env.CausationID.Get(),
	)
	if bkt == nil {
		return nil
	}
	return bkt.Delete([]byte(env.MessageID.Get()))
}

func parseOutboxMessage(
	p []byte,
	causationID ax.MessageID,
) (endpoint.OutboundEnvelope, error) {
	var (
		err error
		m   OutboxMessage
		env endpoint.OutboundEnvelope
	)

	if err = proto.Unmarshal(p, &m); err != nil {
		return env, err
	}

	env.Envelope, err = ax.NewEnvelopeFromProto(m.Envelope)
	if err != nil {
		return env, err
	}

	env.Operation = endpoint.Operation(m.GetOperation())
	env.DestinationEndpoint = m.GetDestinationEndpoint()

	return env, nil
}

func insertOutboxMessage(
	bkt *bolt.Bucket,
	env endpoint.OutboundEnvelope,
) error {
	var err error
	envproto, err := env.Envelope.AsProto()
	if err != nil {
		return err
	}

	m := &OutboxMessage{
		Envelope:            envproto,
		Operation:           int32(env.Operation),
		DestinationEndpoint: env.DestinationEndpoint,
	}

	return boltutil.MarshalProto(
		bkt,
		[]byte(envproto.MessageId),
		m,
	)
}
