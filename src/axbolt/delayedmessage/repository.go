package delayedmessage

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// Repository is a MySQL-backed implementation of Ax's delayedmessage.Repository
// interface.
type Repository struct{}

// DelayedMessageBktName is the name of of the Bolt root bucket where delayed
// messages are stored.
var DelayedMessageBktName = []byte("ax_delayed_message")

// BySendAtBktName is the name of the subbucket in DelayedMessageBktName where
// delayed messages are stored and indexed by SendAt field in message's envelope
var BySendAtBktName = []byte("by_send_at")

// ByIDBktName is the name of the subbucket in DelayedMessageBktName where
// delayed messages are stored and indexed by MessageID field in message's
// envelope
var ByIDBktName = []byte("by_id")

// LoadNextMessage loads the next that is scheduled to be sent.
func (Repository) LoadNextMessage(
	ctx context.Context,
	ds persistence.DataStore,
) (endpoint.OutboundEnvelope, bool, error) {
	db := boltpersistence.ExtractDB(ds)
	tx, err := db.Begin(false)
	if err != nil {
		return endpoint.OutboundEnvelope{}, false, err
	}
	defer tx.Rollback()

	bkt := tx.Bucket(DelayedMessageBktName)
	if bkt == nil {
		return endpoint.OutboundEnvelope{}, false, nil
	}

	if bkt = bkt.Bucket(BySendAtBktName); bkt == nil {
		return endpoint.OutboundEnvelope{}, false, nil
	}

	k, pb := bkt.Cursor().First()
	if k != nil && pb == nil {
		return endpoint.OutboundEnvelope{}, false, nil
	}

	m, err := parseDelayedMessage(pb)
	if err != nil {
		return endpoint.OutboundEnvelope{}, false, err
	}
	return m, true, nil
}

// SaveMessage saves a message to be sent at a later time.
// If does NOT return an error if the message already exists in the repository.
func (Repository) SaveMessage(
	ctx context.Context,
	ptx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {

	tx := boltpersistence.ExtractTx(ptx)
	dmbkt, err := tx.CreateBucketIfNotExists(DelayedMessageBktName)
	if err != nil {
		return err
	}

	bkt, err := dmbkt.CreateBucketIfNotExists(ByIDBktName)
	if err != nil {
		return err
	}

	// return nil if this is a duplicate entry
	if m := bkt.Get([]byte(env.MessageID.Get())); m != nil {
		return nil
	}

	m := &DelayedMessage{
		Id:                  env.MessageID.Get(),
		CausationId:         env.CausationID.Get(),
		CorrelationId:       env.CorrelationID.Get(),
		CreatedAt:           env.CreatedAt.Format(time.RFC3339Nano),
		SendAt:              env.SendAt.Format(time.RFC3339Nano),
		Operation:           int32(env.Operation),
		DestinationEndpoint: env.DestinationEndpoint,
	}
	if m.Message, err = ptypes.MarshalAny(env.Message); err != nil {
		return err
	}

	if err = bkt.Put(
		[]byte(env.MessageID.Get()),
		[]byte(m.SendAt),
	); err != nil {
		return err
	}

	if bkt, err = dmbkt.CreateBucketIfNotExists(BySendAtBktName); err != nil {
		return err
	}

	return boltutil.MarshalProto(bkt, []byte(m.SendAt), m)
}

// MarkAsSent marks a message as sent, removing it from the repository.
func (Repository) MarkAsSent(
	ctx context.Context,
	ptx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {

	tx := boltpersistence.ExtractTx(ptx)
	dmbkt := tx.Bucket(DelayedMessageBktName)
	if dmbkt == nil {
		return nil
	}
	bkt := dmbkt.Bucket(ByIDBktName)
	if bkt == nil {
		return nil
	}
	sa := bkt.Get([]byte(env.MessageID.Get()))
	if sa == nil {
		return nil
	}
	// delete in 'by name' bucket
	if err := bkt.Delete([]byte(env.MessageID.Get())); err != nil {
		return err
	}
	if bkt = dmbkt.Bucket(BySendAtBktName); bkt == nil {
		return nil
	}
	// delete in 'by send at' bucket
	return bkt.Delete(sa)
}

func parseDelayedMessage(
	p []byte,
) (endpoint.OutboundEnvelope, error) {
	var (
		env endpoint.OutboundEnvelope
		dm  DelayedMessage
	)
	err := proto.Unmarshal(p, &dm)
	if err != nil {
		return env, err
	}

	var x ptypes.DynamicAny
	if err = ptypes.UnmarshalAny(dm.Message, &x); err != nil {
		return env, err
	}
	env.Message, _ = x.Message.(ax.Message)

	if err = env.MessageID.Parse(dm.GetId()); err != nil {
		return env, err
	}
	if err = env.CausationID.Parse(dm.GetCausationId()); err != nil {
		return env, err
	}
	if err = env.CorrelationID.Parse(dm.GetCorrelationId()); err != nil {
		return env, err
	}
	env.Operation = endpoint.Operation(dm.GetOperation())
	env.DestinationEndpoint = dm.GetDestinationEndpoint()

	if env.CreatedAt, err = time.Parse(time.RFC3339Nano, dm.GetCreatedAt()); err != nil {
		return env, err
	}

	env.SendAt, err = time.Parse(time.RFC3339Nano, dm.GetSendAt())
	return env, err
}
