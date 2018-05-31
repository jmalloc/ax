package validation

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax"
)

// ValidationError is an error type that contains specifics about message validation.
// Typically Error designates an unrecoverable message state that should
// be retried within either an outbound or inbound message pipeline
type ValidationError struct {
	msg ax.Message
	err error
}

// NewValidationError returns a pointer to a new ValidationError struct
func NewValidationError(err error, msg ax.Message) *ValidationError {
	return &ValidationError{
		msg: msg,
		err: err,
	}
}

// Error returns a string message of an error
func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"validation error: %v",
		e.err,
	)
}
