package validation

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
)

// DefaultValidator is the default message validator
// that performs basic checks on the mesage
type DefaultValidator struct {
}

// Validate validates the message by determining if the message
// implements SelfValidatingMessage interface
func (v *DefaultValidator) Validate(
	ctx context.Context,
	msg ax.Message,
) error {

	// first check if the context is already canceled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// check if message is nil, if it is return a non-recoverable error
	if msg == nil {
		return NewError(errors.New("message cannot be nil"), nil)
	}

	// check if message can perform self-validation
	svmsg, ok := msg.(SelfValidatingMessage)
	if ok {
		return svmsg.Validate()
	}

	return nil
}
