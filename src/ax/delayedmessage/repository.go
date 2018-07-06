package delayedmessage

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Repository is an interface for
type Repository interface {
	// LoadNextMessage loads the next that is scheduled to be sent.
	LoadNextMessage(ctx context.Context, ds persistence.DataStore) (endpoint.OutboundEnvelope, bool, error)

	// SaveMessages saves a message to be sent at a later time.
	SaveMessage(ctx context.Context, tx persistence.Tx, env endpoint.OutboundEnvelope) error
}
