package validation

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax"
)

// Error is an error type that contains specifics about message validation.
// Typically Error designates an unrecoverable message state that cannot
// be retried within either an outbound or inbound message pipeline
type Error struct {
	msg ax.Message
	err error
}

// NewError returns a pointer to a new Error struct
func NewError(err error, msg ax.Message) *Error {
	return &Error{
		msg: msg,
		err: err,
	}
}

// Error returns a string message of an error
func (e *Error) Error() string {
	return fmt.Sprintf(
		"validation error: %v",
		e.err,
	)
}
