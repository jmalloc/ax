package endpoint

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax"
)

// ValidationError is an error type that contains specifics about message validation.
// Typically Error designates an unrecoverable message state that should
// be retried within either an outbound or inbound message pipeline
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

// Error returns a string containing an error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"validation error: %s",
		e.s,
	)
}
