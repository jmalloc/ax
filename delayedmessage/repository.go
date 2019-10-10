package delayedmessage

import (
	"context"

	"github.com/jmalloc/ax/endpoint"
	"github.com/jmalloc/ax/persistence"
)

// Repository is an interface for
type Repository interface {
	// LoadNextMessage loads the next that is scheduled to be sent.
	LoadNextMessage(
		ctx context.Context,
		ds persistence.DataStore,
	) (endpoint.OutboundEnvelope, bool, error)

	// SaveMessage saves a message to be sent at a later time.
	// If does NOT return an error if the message already exists in the repository.
	SaveMessage(
		ctx context.Context,
		tx persistence.Tx,
		env endpoint.OutboundEnvelope,
	) error

	// MarkAsSent marks a message as sent, removing it from the repository.
	MarkAsSent(
		ctx context.Context,
		tx persistence.Tx,
		env endpoint.OutboundEnvelope,
	) error
}
