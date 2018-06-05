package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Validator is an interface for validating messages.
//
// Application-defined validators can be implemented to provide superficial and
// domain validation. Each endpoint has a set of validators that are used to
// validate outgoing messages. Additionally, the validation.InboundStage and
// validation.OutboundStage can be used to perform message validation at any
// point in a pipeline.
type Validator interface {
	// Validate checks if m is valid.

	// It returns a non-nil error if the message is invalid. The meaning of
	// 'valid' in is implementation-defined.
	Validate(ctx context.Context, m ax.Message) error
}

// DefaultValidators is the set of validators used to validate outgoing messages
// if no other set of validators is configured on the endpoint.
var DefaultValidators = []Validator{
	&SelfValidator{},
}

// SelfValidator is one of the default message validators
// that validates the message if it implements SelfValidatingMessage interface.
type SelfValidator struct{}

// Validate validates the message by checking if the message
// is nil and if the message implements SelfValidatingMessage
// interface to call its Validate method.
func (SelfValidator) Validate(
	ctx context.Context,
	msg ax.Message,
) error {

	// check if message can perform self-validation
	if s, ok := msg.(SelfValidatingMessage); ok {
		return s.Validate()
	}

	return nil
}
