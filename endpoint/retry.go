package endpoint

import (
	"math"
	"math/rand"
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
// exponentially until 5-minute maximum delay interval is reached. After that
// the message is retried indefinitely with the maximum delay interval. Every
// single delay produced by this policy is randomized within the following
// range: [ d-(d*02), d+(d*02) ], where d is the calculated delay.
var DefaultRetryPolicy = NewIndefiniteRetryPolicy(3, 1*time.Second, 5*time.Minute, 0.2)

// NewExponentialBackoffPolicy returns a retry policy that allows a fixed number
// of immediate attempts after which retries are delayed exponentially for a
// fixed number of total attempts before the message is rejected.
//
// i is the number of immediate attempts. m is the maximum total attempts, and
// d is a multiplier for the backoff duration.
func NewExponentialBackoffPolicy(i, m uint, d time.Duration) RetryPolicy {
	return func(env InboundEnvelope, _ error) (time.Duration, bool) {
		n := env.AttemptCount

		// Stop retrying if we've reached the maximum number of attempts.
		if n >= m {
			return 0, false
		}

		// If the attempt count is unknown, always retry, but always use the
		// maximum backoff period.
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

// NewIndefiniteRetryPolicy is the policy that regulates indefinite amount of
// retry attempts. Similar to NewExponentialBackoffPolicy, this policy allows a
// fixed number of immediate attempts and calculates the next delay
// exponentially. However, contrary to NewExponentialBackoffPolicy, this policy
// always returns true as the second return value denoting that the message
// should be retried indefinitely.
//
// i is the number of immediate attempts. d is a multiplier for the backoff
// duration. dmax is the maximum backoff duration. randomize is the
// randomization factor used to randomize the returned duration within the
// following range:
//
// [ d-(d*randomize), d+(d*randomize) ]
//
// If d exceeds dmax, dmax value is used to randomize duration.
func NewIndefiniteRetryPolicy(i uint, d, dmax time.Duration, randomize float64) RetryPolicy {
	// lastDelay is the last computed delay returned by this retry policy
	var lastDelay time.Duration
	return func(env InboundEnvelope, _ error) (time.Duration, bool) {
		n := env.AttemptCount

		// If the attempt count is unknown, always retry, but always use the
		// maximum backoff period.
		if n == 0 {
			lastDelay = dmax
			return randomizeDuration(dmax, randomize), true
		}

		// Retry immediately if we haven't yet exhausted the immediate attempt
		// limit.
		if n < i {
			return 0, true
		}

		// lastDelay has already been set to dmax, randomize dmax and return
		if lastDelay == dmax {
			return randomizeDuration(dmax, randomize), true
		}

		// Calculate the next exponential value
		p := math.Pow(
			2,
			float64(n-i), // number of non-immediate attempts
		)

		if delay := time.Duration(p) * d; delay < dmax {
			lastDelay = delay
			return randomizeDuration(delay, randomize), true
		}
		// If we ended up here, delay is already exceeding dmax value.
		// In this case we set lastDelay to dmax to signify that from now on
		// we will keep returning dmax.
		lastDelay = dmax
		return randomizeDuration(dmax, randomize), true
	}
}

// randomizeDuration randomizes d based on the randomization factor. The
// randomized duration is guranteed to be within the range of
// [ d-(d*randomize), d+(d*randomize) ].
func randomizeDuration(d time.Duration, randomize float64) time.Duration {
	f := float64(d)
	delta := f * randomize
	min, max := f-delta, f+delta
	return time.Duration(min + (rand.Float64() * (max - min + 1)))
}
