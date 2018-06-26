package outbox

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/coreos/bbolt"

	"github.com/jmalloc/ax/src/ax/endpoint"

	"github.com/golang/protobuf/proto"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	boltpersistance "github.com/jmalloc/ax/src/axbolt/persistence"
)

// Repository is a Bolt-backed implementation of Ax's outbox.Repository
// interface.
type Repository struct{}

// LoadOutbox loads the unsent outbound messages that were produced when the
// message identified by id was first delivered.
func (Repository) LoadOutbox(
	ctx context.Context,
	ds persistence.DataStore,
	id ax.MessageID,
) ([]endpoint.OutboundEnvelope, bool, error) {
	db := boltpersistance.ExtractDB(ds)
	tx, err := db.Begin(false)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	obBkt := tx.Bucket([]byte("ax_outbox"))
	if obBkt == nil {
		return nil, false, nil
	}

	msgBkt := obBkt.Bucket([]byte(id.Get()))
	if msgBkt == nil {
		return nil, false, nil
	}
	c := msgBkt.Cursor()

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
	tx := boltpersistance.ExtractTx(ptx)
	obBkt, err := tx.CreateBucketIfNotExists([]byte("ax_outbox"))
	if err != nil {
		return err
	}

	var msgBkt *bolt.Bucket
	if msgBkt = obBkt.Bucket([]byte(id.Get())); msgBkt != nil {
		return ErrOutboxExists
	}

	msgBkt, err = obBkt.CreateBucket([]byte(id.Get()))
	if err != nil {
		return err
	}

	for _, env := range envs {
		if err := insertOutboxMessage(msgBkt, env); err != nil {
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
	tx := boltpersistance.ExtractTx(ptx)
	obBkt := tx.Bucket([]byte("ax_outbox"))
	if obBkt == nil {
		return nil
	}
	outMsgBkt := obBkt.Bucket([]byte(env.CausationID.Get()))
	if outMsgBkt == nil {
		return nil
	}

	return outMsgBkt.Delete([]byte(env.MessageID.Get()))
}

func parseOutboxMessage(
	v []byte,
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

	if err = proto.Unmarshal(v, &outmsg); err != nil {
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

	env.Time, err = time.Parse(time.RFC3339Nano, outmsg.GetTime())
	return env, err
}

func insertOutboxMessage(
	bkt *bolt.Bucket,
	env endpoint.OutboundEnvelope,
) error {
	var (
		err error
		pb  []byte
	)
	outmsg := &OutboxMessage{
		Id:                  env.MessageID.Get(),
		CausationId:         env.CausationID.Get(),
		CorrelationId:       env.CorrelationID.Get(),
		Time:                env.Time.Format(time.RFC3339Nano),
		Operation:           int32(env.Operation),
		DestinationEndpoint: env.DestinationEndpoint,
	}

	outmsg.Message, err = ptypes.MarshalAny(env.Message)
	if err != nil {
		return err
	}

	pb, err = proto.Marshal(outmsg)
	if err != nil {
		return err
	}

	return bkt.Put([]byte(outmsg.Id), pb)
}
