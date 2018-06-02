package endpoint

import (
	"github.com/jmalloc/ax/src/ax"
)

// SelfValidatingMessage represents the message that
// can self-validate. A message can implement this
// interface if validation is required both in inbound
// and outbound pipeline. It may or may not return a
// ValidationError to denote an unrecoverable error and
// prevent any further message retries.
type SelfValidatingMessage interface {
	ax.Message
	Validate() error
}
