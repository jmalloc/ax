package observability_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axtest/testmessages"
	"github.com/jmalloc/ax/endpoint"
	. "github.com/jmalloc/ax/observability"
	"github.com/jmalloc/twelf/src/twelf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ensureLoggingObserverIsObserver Observer = &LoggingObserver{}
)

var _ = Describe("LoggingObserver", func() {
	var (
		logger   = &twelf.BufferedLogger{}
		observer = &LoggingObserver{
			Logger: logger,
		}
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
						Message: "= <message-id>  ∵ <causation-id>  ⋲ <correlation-id>  ▼ ↻  axtest.testmessages.Command? ● test command",
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
						Message: "= <message-id>  ∵ <causation-id>  ⋲ <correlation-id>  ▽ ✖  axtest.testmessages.Command? ● <error> ● test command",
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
						Message: "= <message-id>  ∵ <causation-id>  ⋲ <correlation-id>  ▲    axtest.testmessages.Command? ● test command",
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
						Message: "= <message-id>  ∵ <causation-id>  ⋲ <correlation-id>  △ ✖  axtest.testmessages.Command? ● <error> ● test command",
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
