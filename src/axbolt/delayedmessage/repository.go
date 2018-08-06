package delayedmessage

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/marshaling"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// Repository is a Bolt-backed implementation of Ax's delayedmessage.Repository
// interface.
type Repository struct{}

const (
	// DelayedMessageBktName is the name of of the Bolt root bucket where delayed
	// messages are stored.
	DelayedMessageBktName = "ax_delayed_message"

	// BySendAtBktName is the name of the subbucket in DelayedMessageBktName where
	// delayed messages are stored and indexed by SendAt field in message's envelope
	BySendAtBktName = "by_send_at"

	// ByIDBktName is the name of the subbucket in DelayedMessageBktName where
	// delayed messages are stored and indexed by MessageID field in message's
	// envelope
	ByIDBktName = "by_id"
)

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

	bkt := boltutil.GetBkt(
		tx,
		DelayedMessageBktName,
		BySendAtBktName,
	)
	if bkt == nil {
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
	var err error
	tx := boltpersistence.ExtractTx(ptx)
	// return nil if this is a duplicate entry
	if m := boltutil.Get(
		tx,
		env.MessageID.Get(),
		DelayedMessageBktName,
		ByIDBktName,
	); m != nil {
		return nil
	}

	m := &DelayedMessage{
		Id:                  env.MessageID.Get(),
		CausationId:         env.CausationID.Get(),
		CorrelationId:       env.CorrelationID.Get(),
		CreatedAt:           marshaling.MarshalTime(env.CreatedAt),
		SendAt:              marshaling.MarshalTime(env.SendAt),
		Operation:           int32(env.Operation),
		DestinationEndpoint: env.DestinationEndpoint,
	}
	if m.Message, err = ptypes.MarshalAny(env.Message); err != nil {
		return err
	}

	if err = boltutil.Put(
		tx,
		env.MessageID.Get(),
		[]byte(m.SendAt),
		DelayedMessageBktName,
		ByIDBktName,
	); err != nil {
		return err
	}

	bkt, err := boltutil.MakeBkt(
		tx,
		DelayedMessageBktName,
		BySendAtBktName,
	)
	if err != nil {
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
	sa := boltutil.Get(
		tx,
		env.MessageID.Get(),
		DelayedMessageBktName,
		ByIDBktName,
	)
	if sa == nil {
		return nil
	}

	// delete in 'by name' bucket
	if err := boltutil.Delete(
		tx,
		env.MessageID.Get(),
		DelayedMessageBktName,
		ByIDBktName,
	); err != nil {
		return err
	}

	// delete in 'by send at' bucket
	return boltutil.Delete(
		tx,
		string(sa),
		DelayedMessageBktName,
		BySendAtBktName,
	)
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
