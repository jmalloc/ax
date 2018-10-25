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

	BeforeEach(func() {
		logger.Reset()
	})

	Context("inbound messages", func() {
		env := endpoint.InboundEnvelope{
			Envelope: ax.Envelope{
				MessageID:     ax.MustParseMessageID("<message-id>"),
				CausationID:   ax.MustParseMessageID("<causation-id>"),
				CorrelationID: ax.MustParseMessageID("<correlation-id>"),
				Message:       &testmessages.Command{},
			},
			DeliveryID:    endpoint.MustParseDeliveryID("<delivery-id>"),
			DeliveryCount: 3,
		}

		Describe("BeforeInbound", func() {
			It("logs information about the message", func() {
				observer.BeforeInbound(context.Background(), env)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "recv: test command  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id> del:<delivery-id>#3]",
						IsDebug: false,
					},
				))
			})
		})

		Describe("AfterInbound", func() {
			It("logs information about processing errors", func() {
				err := errors.New("<error>")
				observer.AfterInbound(context.Background(), env, err)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "error: test command  <error>  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id> del:<delivery-id>#3]",
						IsDebug: false,
					},
				))
			})

			It("does not log if no error occurred", func() {
				observer.AfterInbound(context.Background(), env, nil)

				Expect(logger.Messages()).To(BeEmpty())
			})
		})
	})

	Context("outbound messages", func() {
		env := endpoint.OutboundEnvelope{
			Envelope: ax.Envelope{
				MessageID:     ax.MustParseMessageID("<message-id>"),
				CausationID:   ax.MustParseMessageID("<causation-id>"),
				CorrelationID: ax.MustParseMessageID("<correlation-id>"),
				Message:       &testmessages.Command{},
			},
		}

		Describe("AfterOutbound", func() {
			It("logs information about the message", func() {
				observer.AfterOutbound(context.Background(), env, nil)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "send: test command  [axtest.testmessages.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>]",
						IsDebug: false,
					},
				))
			})

			It("does not log if an error occurred", func() {
				err := errors.New("<error>")
				observer.AfterOutbound(context.Background(), env, err)

				Expect(logger.Messages()).To(BeEmpty())
			})
		})
	})
})
