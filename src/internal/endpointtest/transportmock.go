// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package endpointtest

import (
	"context"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"sync"
)

var (
	lockTransportMockInitialize sync.RWMutex
	lockTransportMockReceive    sync.RWMutex
	lockTransportMockSend       sync.RWMutex
	lockTransportMockSubscribe  sync.RWMutex
)

// TransportMock is a mock implementation of Transport.
//
//     func TestSomethingThatUsesTransport(t *testing.T) {
//
//         // make and configure a mocked Transport
//         mockedTransport := &TransportMock{
//             InitializeFunc: func(ctx context.Context, ep string) error {
// 	               panic("TODO: mock out the Initialize method")
//             },
//             ReceiveFunc: func(ctx context.Context) (endpoint.InboundEnvelope, endpoint.Acknowledger, error) {
// 	               panic("TODO: mock out the Receive method")
//             },
//             SendFunc: func(ctx context.Context, env endpoint.OutboundEnvelope) error {
// 	               panic("TODO: mock out the Send method")
//             },
//             SubscribeFunc: func(ctx context.Context, op endpoint.Operation, mt ax.MessageTypeSet) error {
// 	               panic("TODO: mock out the Subscribe method")
//             },
//         }
//
//         // TODO: use mockedTransport in code that requires Transport
//         //       and then make assertions.
//
//     }
type TransportMock struct {
	// InitializeFunc mocks the Initialize method.
	InitializeFunc func(ctx context.Context, ep string) error

	// ReceiveFunc mocks the Receive method.
	ReceiveFunc func(ctx context.Context) (endpoint.InboundEnvelope, endpoint.Acknowledger, error)

	// SendFunc mocks the Send method.
	SendFunc func(ctx context.Context, env endpoint.OutboundEnvelope) error

	// SubscribeFunc mocks the Subscribe method.
	SubscribeFunc func(ctx context.Context, op endpoint.Operation, mt ax.MessageTypeSet) error

	// calls tracks calls to the methods.
	calls struct {
		// Initialize holds details about calls to the Initialize method.
		Initialize []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Ep is the ep argument value.
			Ep string
		}
		// Receive holds details about calls to the Receive method.
		Receive []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Send holds details about calls to the Send method.
		Send []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Env is the env argument value.
			Env endpoint.OutboundEnvelope
		}
		// Subscribe holds details about calls to the Subscribe method.
		Subscribe []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Op is the op argument value.
			Op endpoint.Operation
			// Mt is the mt argument value.
			Mt ax.MessageTypeSet
		}
	}
}

// Initialize calls InitializeFunc.
func (mock *TransportMock) Initialize(ctx context.Context, ep string) error {
	if mock.InitializeFunc == nil {
		panic("moq: TransportMock.InitializeFunc is nil but Transport.Initialize was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Ep  string
	}{
		Ctx: ctx,
		Ep:  ep,
	}
	lockTransportMockInitialize.Lock()
	mock.calls.Initialize = append(mock.calls.Initialize, callInfo)
	lockTransportMockInitialize.Unlock()
	return mock.InitializeFunc(ctx, ep)
}

// InitializeCalls gets all the calls that were made to Initialize.
// Check the length with:
//     len(mockedTransport.InitializeCalls())
func (mock *TransportMock) InitializeCalls() []struct {
	Ctx context.Context
	Ep  string
} {
	var calls []struct {
		Ctx context.Context
		Ep  string
	}
	lockTransportMockInitialize.RLock()
	calls = mock.calls.Initialize
	lockTransportMockInitialize.RUnlock()
	return calls
}

// Receive calls ReceiveFunc.
func (mock *TransportMock) Receive(ctx context.Context) (endpoint.InboundEnvelope, endpoint.Acknowledger, error) {
	if mock.ReceiveFunc == nil {
		panic("moq: TransportMock.ReceiveFunc is nil but Transport.Receive was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	lockTransportMockReceive.Lock()
	mock.calls.Receive = append(mock.calls.Receive, callInfo)
	lockTransportMockReceive.Unlock()
	return mock.ReceiveFunc(ctx)
}

// ReceiveCalls gets all the calls that were made to Receive.
// Check the length with:
//     len(mockedTransport.ReceiveCalls())
func (mock *TransportMock) ReceiveCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	lockTransportMockReceive.RLock()
	calls = mock.calls.Receive
	lockTransportMockReceive.RUnlock()
	return calls
}

// Send calls SendFunc.
func (mock *TransportMock) Send(ctx context.Context, env endpoint.OutboundEnvelope) error {
	if mock.SendFunc == nil {
		panic("moq: TransportMock.SendFunc is nil but Transport.Send was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
	}{
		Ctx: ctx,
		Env: env,
	}
	lockTransportMockSend.Lock()
	mock.calls.Send = append(mock.calls.Send, callInfo)
	lockTransportMockSend.Unlock()
	return mock.SendFunc(ctx, env)
}

// SendCalls gets all the calls that were made to Send.
// Check the length with:
//     len(mockedTransport.SendCalls())
func (mock *TransportMock) SendCalls() []struct {
	Ctx context.Context
	Env endpoint.OutboundEnvelope
} {
	var calls []struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
	}
	lockTransportMockSend.RLock()
	calls = mock.calls.Send
	lockTransportMockSend.RUnlock()
	return calls
}

// Subscribe calls SubscribeFunc.
func (mock *TransportMock) Subscribe(ctx context.Context, op endpoint.Operation, mt ax.MessageTypeSet) error {
	if mock.SubscribeFunc == nil {
		panic("moq: TransportMock.SubscribeFunc is nil but Transport.Subscribe was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Op  endpoint.Operation
		Mt  ax.MessageTypeSet
	}{
		Ctx: ctx,
		Op:  op,
		Mt:  mt,
	}
	lockTransportMockSubscribe.Lock()
	mock.calls.Subscribe = append(mock.calls.Subscribe, callInfo)
	lockTransportMockSubscribe.Unlock()
	return mock.SubscribeFunc(ctx, op, mt)
}

// SubscribeCalls gets all the calls that were made to Subscribe.
// Check the length with:
//     len(mockedTransport.SubscribeCalls())
func (mock *TransportMock) SubscribeCalls() []struct {
	Ctx context.Context
	Op  endpoint.Operation
	Mt  ax.MessageTypeSet
} {
	var calls []struct {
		Ctx context.Context
		Op  endpoint.Operation
		Mt  ax.MessageTypeSet
	}
	lockTransportMockSubscribe.RLock()
	calls = mock.calls.Subscribe
	lockTransportMockSubscribe.RUnlock()
	return calls
}