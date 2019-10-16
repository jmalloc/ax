package endpoint

import (
	"math"
	"time"
)

// RetryPolicy is a function responsible for determining whether or not a
// message should be retried.
//
// It returns the delay that should occur before retrying, and a bool indicating
// whether or not the message should be retried at all.
type RetryPolicy func(InboundEnvelope, error) (time.Duration, bool)

// DefaultRetryPolicy is the default RetryPolicy.
//
// It allows for 3 immediate attempts, after which each attempt is delayed
// exponentially, for a maximum of 10 attempts before the message is rejected.
var DefaultRetryPolicy = NewExponentialBackoffPolicy(
	3,
	0, // retry forever
	1*time.Second,
	1*time.Hour,
)

// NewExponentialBackoffPolicy returns a retry policy that allows a fixed number
// of immediate attempts after which retries are delayed exponentially for a
// until some maximum delay is reached.
//
// Optionally, the message can be rejected after some fixed number of retries.
//
// ir is the number of immediate attempts. mr is the maximum total attempts
// before rejecting the message. If mr is zero, the message is retried
// indefinitely.
//
// bt is a "base" delay between retries. It is used as a multplier for the
// backoff duration. mt is the maximum delay between retries.
func NewExponentialBackoffPolicy(
	ir, mr uint,
	bt, mt time.Duration,
) RetryPolicy {
	return func(env InboundEnvelope, _ error) (time.Duration, bool) {
		n := env.AttemptCount

		// If the attempt count is unknown, always retry, but always use the
		// maximium backoff period.
		if n == 0 {
			return mt, true
		}

		// Stop retrying if we've reached the maximum number of attempts.
		if mr != 0 && n >= mr {
			return 0, false
		}

		// Retry immediately if we haven't yet exhausted the immediate attempt limit.
		if n < ir {
			return 0, true
		}

		// Otherwise, backoff exponentially.
		p := math.Pow(
			2,
			float64(n-ir), // number of non-immediate attempts
		)

		// If the multiplier itself would overflow a time duration, use the cap.
		if p > math.MaxInt64 {
			return mt, true
		}

		// Otherwise; cap the delay at the maximum.
		d := time.Duration(p) * bt
		if d > mt {
			return mt, true
		}

		return d, true
	}
}
