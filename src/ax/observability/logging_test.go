package observability_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	"github.com/jmalloc/twelf/src/twelf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ensureLoggingObserverIsBeforeInboundObserver BeforeInboundObserver = &LoggingObserver{}
	ensureLoggingObserverIsAfterInboundObserver  AfterInboundObserver  = &LoggingObserver{}
	ensureLoggingObserverIsAfterOutboundObserver AfterOutboundObserver = &LoggingObserver{}
)

var _ = Describe("Logger", func() {
	var (
		logger   = &twelf.BufferedLogger{}
		observer = &LoggingObserver{Logger: logger}
	)

	in := endpoint.InboundEnvelope{
		Envelope: ax.Envelope{
			MessageID:     ax.MustParseMessageID("<message-id>"),
			CausationID:   ax.MustParseMessageID("<causation-id>"),
			CorrelationID: ax.MustParseMessageID("<correlation-id>"),
			Message:       &testmessages.Command{},
		},
		AttemptID:    endpoint.MustParseAttemptID("<attempt-id>"),
		AttemptCount: 3,
	}

	out := endpoint.OutboundEnvelope{
		Envelope: ax.Envelope{
			MessageID:     ax.MustParseMessageID("<message-id>"),
			CausationID:   ax.MustParseMessageID("<causation-id>"),
			CorrelationID: ax.MustParseMessageID("<correlation-id>"),
			Message:       &testmessages.Command{},
		},
	}

	BeforeEach(func() {
		logger.Reset()
	})

	Context("inbound messages", func() {
		Describe("BeforeInbound", func() {
			It("logs information about the message", func() {
				ctx := endpoint.WithEnvelope(context.Background(), in)

				observer.BeforeInbound(ctx, in)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "▼   test command  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>] [attempt:<attempt-id> #3]",
						IsDebug: false,
					},
				))
			})
		})

		Describe("AfterInbound", func() {
			It("logs information about errors", func() {
				ctx := endpoint.WithEnvelope(context.Background(), in)
				err := errors.New("<error>")

				observer.AfterInbound(ctx, in, err)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "▽ ✘ test command ∎ <error>  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>] [attempt:<attempt-id> #3]",
						IsDebug: false,
					},
				))
			})

			It("does not log if no error occurred", func() {
				observer.AfterInbound(context.Background(), in, nil)

				Expect(logger.Messages()).To(BeEmpty())
			})
		})
	})

	Context("outbound messages", func() {
		Describe("BeforeOutbound", func() {
			It("logs information about the message", func() {
				observer.BeforeOutbound(context.Background(), out)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "▲   test command  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>]",
						IsDebug: false,
					},
				))
			})

			It("includes information about the attempt if the inbound envelope is in the context", func() {
				ctx := endpoint.WithEnvelope(context.Background(), in)

				observer.BeforeOutbound(ctx, out)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "▲   test command  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>] [attempt:<attempt-id> #3]",
						IsDebug: false,
					},
				))
			})
		})

		Describe("AfterOutbound", func() {
			It("logs information about errors", func() {
				err := errors.New("<error>")
				observer.AfterOutbound(context.Background(), out, err)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "△ ✘ test command ∎ <error>  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>]",
						IsDebug: false,
					},
				))
			})

			It("includes information about the attempt if the inbound envelope is in the context", func() {
				ctx := endpoint.WithEnvelope(context.Background(), in)
				err := errors.New("<error>")

				observer.AfterOutbound(ctx, out, err)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "△ ✘ test command ∎ <error>  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>] [attempt:<attempt-id> #3]",
						IsDebug: false,
					},
				))
			})

			It("does not log if no error occurred", func() {
				observer.AfterOutbound(context.Background(), out, nil)

				Expect(logger.Messages()).To(BeEmpty())
			})
		})
	})
})
