package routing_test

import (
	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/routing"
	"github.com/jmalloc/ax/src/internal/messagetest"
	"github.com/jmalloc/ax/src/internal/routingtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HandlerTable", func() {
	var table HandlerTable

	Describe("NewHandlerTable", func() {
		It("returns an error when multiple handlers accept the same command", func() {
			h1 := &routingtest.MessageHandlerMock{
				MessageTypesFunc: func() ax.MessageTypeSet {
					return ax.TypesOf(&messagetest.Command{})
				},
			}

			h2 := &routingtest.MessageHandlerMock{
				MessageTypesFunc: func() ax.MessageTypeSet {
					return ax.TypesOf(&messagetest.Command{})
				},
			}

			_, err := NewHandlerTable(h1, h2)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Lookup", func() {
		h1 := &routingtest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&messagetest.Command{})
			},
		}

		h2 := &routingtest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&messagetest.Event{})
			},
		}

		h3 := &routingtest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&messagetest.Event{})
			},
		}

		BeforeEach(func() {
			t, err := NewHandlerTable(h1, h2, h3)
			Expect(err).ShouldNot(HaveOccurred())
			table = t
		})

		It("returns all of the handlers that handle the given message type", func() {
			mt := ax.TypeOf(&messagetest.Event{})
			h := table.Lookup(mt)

			Expect(h).To(ConsistOf(h2, h3))
		})
	})
})
