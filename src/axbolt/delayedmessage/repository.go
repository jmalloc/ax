package delayedmessage

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// Repository is a Bolt-backed implementation of Ax's delayedmessage.Repository
// interface.
type Repository struct{}

const (
	// delayedMessageBktName is the name of the Bolt root bucket where delayed
	// messages are stored.
	delayedMessageBktName = "ax_delayed_message"

	// bySendAtBktName is the name of a subbucket in delayedMessageBktName where
	// delayed messages are stored and indexed by SendAt field in message's
	// envelope
	bySendAtBktName = "by_send_at"

	// byIDBktName is the name of a subbucket in delayedMessageBktName where
	// delayed messages are stored and indexed by MessageID field in message's
	// envelope
	byIDBktName = "by_id"
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
		delayedMessageBktName,
		bySendAtBktName,
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
		delayedMessageBktName,
		byIDBktName,
	); m != nil {
		return nil
	}

	envproto, err := env.AsProto()
	if err != nil {
		return err
	}

	m := &DelayedMessage{
		Envelope:            envproto,
		Operation:           int32(env.Operation),
		DestinationEndpoint: env.DestinationEndpoint,
	}

	if err = boltutil.Put(
		tx,
		env.MessageID.Get(),
		[]byte(ptypes.TimestampString(envproto.SendAt)),
		delayedMessageBktName,
		byIDBktName,
	); err != nil {
		return err
	}

	return boltutil.PutProto(
		tx,
		ptypes.TimestampString(envproto.SendAt),
		m,
		delayedMessageBktName,
		bySendAtBktName,
	)
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
		delayedMessageBktName,
		byIDBktName,
	)
	if sa == nil {
		return nil
	}

	// delete in 'by name' bucket
	if err := boltutil.Delete(
		tx,
		env.MessageID.Get(),
		delayedMessageBktName,
		byIDBktName,
	); err != nil {
		return err
	}

	// delete in 'by send at' bucket
	return boltutil.Delete(
		tx,
		string(sa),
		delayedMessageBktName,
		bySendAtBktName,
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

	env.Envelope, err = ax.NewEnvelopeFromProto(dm.GetEnvelope())
	if err != nil {
		return env, err
	}
	env.Operation = endpoint.Operation(dm.GetOperation())
	env.DestinationEndpoint = dm.GetDestinationEndpoint()

	return env, err
}
