package endpoint

import (
	"github.com/jmalloc/ax/src/ax"
)

// SelfValidatingMessage is a message that can perform its own superficial
// validation.
type SelfValidatingMessage interface {
	ax.Message
	// Validate returns a non-nil error if the message is invalid. It is up to
	// message implementation to check validity criteria.
	//
	// This method is invoked by SelfValidator that is one the default
	// validators to verify the message.
	Validate() error
}
