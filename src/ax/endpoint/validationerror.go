package endpoint

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax"
)

// ValidationError is an error type that contains specifics about message validation.
// Typically, ValidationError designates an unrecoverable message state that should not
// be retried within either an outbound or inbound message pipelines.
type ValidationError struct {
	InvalidMsg ax.Message
	s          string
}

// NewValidationError returns a pointer to a new ValidationError struct
func NewValidationError(s string, msg ax.Message) *ValidationError {
	return &ValidationError{
		s:          s,
		InvalidMsg: msg,
	}
}

// Error returns a string containing a validation error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"validation error: %s",
		e.s,
	)
}
