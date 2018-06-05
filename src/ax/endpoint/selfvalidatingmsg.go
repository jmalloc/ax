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

// SelfValidatingCommand is a command that can perform its own superficial
// validation.
type SelfValidatingCommand interface {
	ax.Command
	// Validate returns a non-nil error if the command is invalid. It is up to
	// command implementation to check validity criteria.
	//
	// This method is invoked by SelfValidator that is one the default
	// validators to verify the message.
	Validate() error
}

// SelfValidatingEvent is an event that can perform its own superficial
// validation.
type SelfValidatingEvent interface {
	ax.Event
	// Validate returns a non-nil error if an event is invalid. It is up to
	// event implementation to check validity criteria.
	//
	// This method is invoked by SelfValidator that is one the default
	// validators to verify the message.
	Validate() error
}
