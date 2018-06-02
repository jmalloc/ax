package endpoint

import (
	"github.com/jmalloc/ax/src/ax"
)

// SelfValidatingMessage represents the message that
// can perform a self validation. A message can implement
// this interface if validation is required both in inbound
// and outbound pipeline.
type SelfValidatingMessage interface {
	ax.Message
	Validate() error
}
