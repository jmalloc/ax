package outbox

import (
	"context"
	"errors"
	"time"

	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"

	"github.com/golang/protobuf/ptypes"

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

// OutboxBktName is the name of of the Bolt root bucket where all
// outbox-specific data is stored.
var OutboxBktName = []byte("ax_outbox")

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

	bkt := tx.Bucket(OutboxBktName)
	if bkt == nil {
		return nil, false, nil
	}

	if bkt = bkt.Bucket([]byte(id.Get())); bkt == nil {
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
	bkt, err := tx.CreateBucketIfNotExists(OutboxBktName)
	if err != nil {
		return err
	}

	if b := bkt.Bucket([]byte(id.Get())); b != nil {
		return ErrOutboxExists
	}

	bkt, err = bkt.CreateBucket([]byte(id.Get()))
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
	bkt := tx.Bucket(OutboxBktName)
	if bkt == nil {
		return nil
	}
	if bkt = bkt.Bucket([]byte(env.CausationID.Get())); bkt == nil {
		return nil
	}

	return bkt.Delete([]byte(env.MessageID.Get()))
}

func parseOutboxMessage(
	p []byte,
	causationID ax.MessageID,
) (endpoint.OutboundEnvelope, error) {
	var (
		err    error
		outmsg OutboxMessage
	)
	env := endpoint.OutboundEnvelope{
		Envelope: ax.Envelope{
			CausationID: causationID,
		},
	}

	if err = proto.Unmarshal(p, &outmsg); err != nil {
		return env, err
	}

	var x ptypes.DynamicAny
	if err = ptypes.UnmarshalAny(outmsg.Message, &x); err != nil {
		return env, err
	}
	env.Message, _ = x.Message.(ax.Message)

	if err = env.MessageID.Parse(outmsg.GetId()); err != nil {
		return env, err
	}
	if err = env.CorrelationID.Parse(outmsg.GetCorrelationId()); err != nil {
		return env, err
	}
	env.Operation = endpoint.Operation(outmsg.GetOperation())
	env.DestinationEndpoint = outmsg.GetDestinationEndpoint()

	if env.CreatedAt, err = time.Parse(time.RFC3339Nano, outmsg.GetCreatedAt()); err != nil {
		return env, err
	}

	env.SendAt, err = time.Parse(time.RFC3339Nano, outmsg.GetSendAt())

	return env, err
}

func insertOutboxMessage(
	bkt *bolt.Bucket,
	env endpoint.OutboundEnvelope,
) error {
	var err error

	m := &OutboxMessage{
		Id:                  env.MessageID.Get(),
		CausationId:         env.CausationID.Get(),
		CorrelationId:       env.CorrelationID.Get(),
		CreatedAt:           env.CreatedAt.Format(time.RFC3339Nano),
		SendAt:              env.SendAt.Format(time.RFC3339Nano),
		Operation:           int32(env.Operation),
		DestinationEndpoint: env.DestinationEndpoint,
	}

	m.Message, err = ptypes.MarshalAny(env.Message)
	if err != nil {
		return err
	}

	return boltutil.MarshalProto(bkt, []byte(m.GetId()), m)
}
