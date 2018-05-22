package validation

import "github.com/jmalloc/ax/src/ax"

// ValidationError is an error type that contains
// specifics about superficial message validation
type ValidationError struct {
	msg ax.Message
	s   string
}

// NewValidationError returns a pointer to a new ValidationError struct
func NewValidationError(s string, msg ax.Message) *ValidationError {
	return &ValidationError{
		s:   s,
		msg: msg,
	}
}

// Error returns a string message of an error
func (ve *ValidationError) Error() string {
	return ve.s
}
