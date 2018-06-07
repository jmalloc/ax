package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Validator is an interface for validating messages.
//
// Application-defined validators can be implemented to provide superficial and
// domain validation. Each endpoint has a set of validators that are used to
// validate outgoing messages. Additionally, the validation.InboundRejecter and
// validation.OutboundRejecter can be used to perform message validation at any
// point in a pipeline.
type Validator interface {
	// Validate checks if m is valid.
	//
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

// Validate validates m if it implements SelfValidatingMessage. It returns the
// error returned by m.Validate(). If m does not implement SelfValidatingMessage
// then no validation is performed and nil is returned.
func (SelfValidator) Validate(
	ctx context.Context,
	m ax.Message,
) error {

	// check if message can perform self-validation
	if s, ok := m.(SelfValidatingMessage); ok {
		return s.Validate()
	}

	return nil
}
