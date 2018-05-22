package validation

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax"
)

// ValidationError is an error type that contains
// specifics about superficial message validation
type ValidationError struct {
	env   ax.Envelope
	cause error
}

// NewValidationError returns a pointer to a new ValidationError struct
func NewValidationError(err error, env ax.Envelope) *ValidationError {
	return &ValidationError{
		env:   env,
		cause: err,
	}
}

// Error returns a string message of an error
func (ve *ValidationError) Error() string {
	return fmt.Sprintf(
		"validation error, cause: %v",
		ve.cause,
	)
}
