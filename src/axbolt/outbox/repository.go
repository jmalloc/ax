package outbox

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/coreos/bbolt"

	"github.com/jmalloc/ax/src/ax/endpoint"

	"github.com/golang/protobuf/proto"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	boltpersistance "github.com/jmalloc/ax/src/axbolt/persistence"
)

// Repository is a MySQL-backed implementation of Ax's outbox.Repository
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

	outboxBkt := tx.Bucket([]byte("ax_outbox"))
	if outboxBkt == nil {
		return nil, false, nil
	}

	outMsgBkt := outboxBkt.Bucket([]byte(id.String()))
	if outMsgBkt == nil {
		return nil, false, nil
	}
	c := outMsgBkt.Cursor()

	var envelopes []endpoint.OutboundEnvelope
	for k, v := c.First(); k != nil; k, v = c.Next() {
		env, err := parseOutboxMessage(v, id)
		if err != nil {
			return nil, false, nil
		}
		envelopes = append(envelopes, env)
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}

	return envelopes, true, nil
}

// SaveOutbox saves a set of unsent outbound messages that were produced
// when the message identified by id was delivered.
func (Repository) SaveOutbox(
	ctx context.Context,
	ptx persistence.Tx,
	id ax.MessageID,
	envs []endpoint.OutboundEnvelope,
) error {
	if len(envs) == 0 {
		return nil
	}

	tx := boltpersistance.ExtractTx(ptx)
	outboxBkt, err := tx.CreateBucketIfNotExists([]byte("ax_outbox"))
	if err != nil {
		return err
	}

	outMsgBkt, err := outboxBkt.CreateBucketIfNotExists([]byte(id.String()))
	if err != nil {
		return err
	}

	for _, env := range envs {
		if err := insertOutboxMessage(outMsgBkt, env); err != nil {
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
	outboxBkt := tx.Bucket([]byte("ax_outbox"))
	if outboxBkt == nil {
		return nil
	}
	outMsgBkt := outboxBkt.Bucket([]byte(env.CausationID.String()))
	if outMsgBkt == nil {
		return nil
	}

	return outMsgBkt.Delete([]byte(env.CausationID.String()))
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

	env.Operation = endpoint.Operation(outmsg.GetOperation())
	env.DestinationEndpoint = outmsg.GetDestinationEndpoint()

	if err = ptypes.UnmarshalAny(outmsg.Message, env.Message); err != nil {
		return env, err
	}

	if err = env.MessageID.Parse(outmsg.GetId()); err != nil {
		return env, err
	}

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
		Id:                  env.MessageID.String(),
		CausationId:         env.CausationID.String(),
		CorrelationId:       env.CorrelationID.String(),
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
