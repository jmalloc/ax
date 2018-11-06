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
	ensureLoggingObserverIsInboundObserver  InboundObserver  = &LoggingObserver{}
	ensureLoggingObserverIsOutboundObserver OutboundObserver = &LoggingObserver{}
)

var _ = Describe("Logger", func() {
	var (
		logger   = &twelf.BufferedLogger{}
		observer = &LoggingObserver{}
		ep       = &endpoint.Endpoint{
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
		BeforeEach(func() {
			if err := observer.InitializeInbound(context.Background(), ep); err != nil {
				panic(err)
			}
		})

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
		BeforeEach(func() {
			if err := observer.InitializeOutbound(context.Background(), ep); err != nil {
				panic(err)
			}
		})

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
