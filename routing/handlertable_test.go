package routing_test

import (
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axtest/mocks"
	"github.com/jmalloc/ax/axtest/testmessages"
	. "github.com/jmalloc/ax/routing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HandlerTable", func() {
	var table HandlerTable

	Describe("NewHandlerTable", func() {
		It("returns an error when multiple handlers accept the same command", func() {
			h1 := &mocks.MessageHandlerMock{
				MessageTypesFunc: func() ax.MessageTypeSet {
					return ax.TypesOf(&testmessages.Command{})
				},
			}

			h2 := &mocks.MessageHandlerMock{
				MessageTypesFunc: func() ax.MessageTypeSet {
					return ax.TypesOf(&testmessages.Command{})
				},
			}

			_, err := NewHandlerTable(h1, h2)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Lookup", func() {
		h1 := &mocks.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&testmessages.Command{})
			},
		}

		h2 := &mocks.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&testmessages.Event{})
			},
		}

		h3 := &mocks.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&testmessages.Event{})
			},
		}

		BeforeEach(func() {
			t, err := NewHandlerTable(h1, h2, h3)
			Expect(err).ShouldNot(HaveOccurred())
			table = t
		})

		It("returns all of the handlers that handle the given message type", func() {
			mt := ax.TypeOf(&testmessages.Event{})
			h := table.Lookup(mt)

			Expect(h).To(ConsistOf(h2, h3))
		})
	})
})
