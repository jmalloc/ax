package endpoint

import (
	"github.com/jmalloc/ax/src/ax/bus"
)

// RetryPolicy returns true if the message should be retried.
type RetryPolicy func(bus.InboundEnvelope, error) bool

// DefaultRetryPolicy is a RetryPolicy that rejects a message after it has been
// attempted three (3) times.
func DefaultRetryPolicy(env bus.InboundEnvelope, _ error) bool {
	return env.DeliveryCount < 3
}
