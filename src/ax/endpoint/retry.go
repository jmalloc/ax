package endpoint

import (
	"math"
	"time"
)

// RetryPolicy returns is a function responsible for determining whether or not
// a message should be retried.
//
// It returns the delay that occur before retrying, and a bool indicating
// whether or not the message should be retried at all.
type RetryPolicy func(InboundEnvelope, error) (time.Duration, bool)

// DefaultRetryPolicy is the default RetryPolicy.
//
// It allows for 3 immediate attempts, after which each attempt is delayed
// exponentially, for a maximum of 10 attempts before the message is rejected.
var DefaultRetryPolicy = NewExponentialBackoffPolicy(3, 10, 1*time.Second)

// NewExponentialBackoffPolicy returns a retry policy that allows a fixed number
// of immediate attempts after which retries are delayed exponentially for a
// fixed number of total attempts before the message is rejected.
//
// i is the number of immediate attempts. m is the maximum total attempts, and
// d is a multplier for the backoff duration.
func NewExponentialBackoffPolicy(i, m uint, d time.Duration) RetryPolicy {
	return func(env InboundEnvelope, _ error) (time.Duration, bool) {
		n := env.DeliveryCount

		// Stop retrying if we've reached the maximum number of attempts.
		if n >= m {
			return 0, false
		}

		// If the delivery count is unknown, always retry, but always use the
		// maximium backoff period.
		if n == 0 {
			n = m
		}

		// Retry immediately if we haven't yet exhausted the immediate attempt limit.
		if n < i {
			return 0, true
		}

		// Otherwise, backoff exponentially.
		p := math.Pow(
			2,
			float64(n-i), // number of non-immediate attempts
		)

		return time.Duration(p) * d, true
	}
}
