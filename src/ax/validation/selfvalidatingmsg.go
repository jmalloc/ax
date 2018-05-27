package validation

import (
	"github.com/jmalloc/ax/src/ax"
)

// SelfValidatingMessage represents the message that
// can perform a self validation. The implementation of
// this interface is typically checked within validators.
type SelfValidatingMessage interface {
	ax.Message
	Validate() error
}
