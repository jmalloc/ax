package endpoint_test

import (
	"time"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/jmalloc/ax/endpoint"
	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable(
	"NewExponentialBackoffPolicy",
	func(
		attempts int,
		ir, mr int,
		retry bool,
		delay time.Duration,
	) {
		policy := NewExponentialBackoffPolicy(
			uint(ir),
			uint(mr),
			3*time.Second,
			10*time.Minute,
		)

		env := InboundEnvelope{
			AttemptCount: uint(attempts),
		}

		d, ok := policy(env, nil)
		Expect(ok).To(Equal(retry))

		if ok {
			Expect(d).To(Equal(delay))
		}
	},
	Entry(
		"retries immediately",
		2,             // attempt count
		3,             // immediate reties
		0,             // max retries
		true,          // expect retry
		0*time.Second, // expected delay
	),
	Entry(
		"does not retry immediately if the attempt count reaches ir",
		3,             // attempt count
		3,             // immediate reties
		0,             // max retries
		true,          // expect retry
		3*time.Second, // expected delay
	),
	Entry(
		"stops retrying if the attempt count reaches mr",
		5,             // attempt count
		3,             // immediate reties
		5,             // max retries
		false,         // expect retry
		0*time.Second, // expected delay
	),
	Entry(
		"interpolates the delay exponentially",
		5,              // attempt count
		3,              // immediate reties
		0,              // max retries
		true,           // expect retry
		12*time.Second, // expected delay
	),
	Entry(
		"caps the delay at mt",
		50,             // attempt count
		3,              // immediate reties
		0,              // max retries
		true,           // expect retry
		10*time.Minute, // expected delay
	),
	Entry(
		"caps the delay at mt if the float exponent overflows the duration type",
		10000,          // attempt count
		3,              // immediate reties
		0,              // max retries
		true,           // expect retry
		10*time.Minute, // expected delay
	),
	Entry(
		"always uses mt if the attempt count is unknown ",
		0,              // attempt count
		3,              // immediate reties
		0,              // max retries
		true,           // expect retry
		10*time.Minute, // expected delay
	),
)
