package observability_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	. "github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/internal/messagetest"
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
		env := bus.InboundEnvelope{
			Envelope: ax.Envelope{
				Message: &messagetest.Command{},
			},
		}

		env.MessageID.MustParse("<message-id>")
		env.CausationID.MustParse("<causation-id>")
		env.CorrelationID.MustParse("<correlation-id>")

		Describe("BeforeInbound", func() {
			It("logs information about the message", func() {
				observer.BeforeInbound(context.Background(), env)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "recv: test command  [ax.internal.messagetest.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>]",
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
						Message: "error: test command  <error>  [ax.internal.messagetest.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>]",
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
		env := bus.OutboundEnvelope{
			Envelope: ax.Envelope{
				Message: &messagetest.Command{},
			},
		}

		env.MessageID.MustParse("<message-id>")
		env.CausationID.MustParse("<causation-id>")
		env.CorrelationID.MustParse("<correlation-id>")

		Describe("AfterOutbound", func() {
			It("logs information about the message", func() {
				observer.AfterOutbound(context.Background(), env, nil)

				Expect(logger.Messages()).To(ConsistOf(
					twelf.BufferedLogMessage{
						Message: "send: test command  [ax.internal.messagetest.Command? msg:<message-id> cause:<causation-id> corr:<correlation-id>]",
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
