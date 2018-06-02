package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Validator wraps methods required to
// perform message validation
type Validator interface {
	Validate(ctx context.Context, msg ax.Message) error
}

// DefaultValidators is the array of validators used by default
// when sending a message
var DefaultValidators = []Validator{
	&BasicValidator{},
}

// BasicValidator is one of  the default message validators
// that performs basic checks on the message such as
// if a message is nil, etc.
type BasicValidator struct{}

// Validate validates the message by determining if the message
// implements SelfValidatingMessage interface
func (v *BasicValidator) Validate(
	ctx context.Context,
	msg ax.Message,
) error {

	// check if message is nil, if it is return a non-recoverable error
	if msg == nil {
		return NewValidationError("message cannot be nil", nil)
	}

	// check if message can perform self-validation
	svmsg, ok := msg.(SelfValidatingMessage)
	if ok {
		return svmsg.Validate()
	}

	return nil
}
