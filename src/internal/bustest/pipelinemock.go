// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package bustest

import (
	"context"
	"github.com/jmalloc/ax/src/ax/bus"
	"sync"
)

var (
	lockInboundPipelineMockAccept     sync.RWMutex
	lockInboundPipelineMockInitialize sync.RWMutex
)

// InboundPipelineMock is a mock implementation of InboundPipeline.
//
//     func TestSomethingThatUsesInboundPipeline(t *testing.T) {
//
//         // make and configure a mocked InboundPipeline
//         mockedInboundPipeline := &InboundPipelineMock{
//             AcceptFunc: func(in1 context.Context, in2 bus.MessageSink, in3 bus.InboundEnvelope) error {
// 	               panic("TODO: mock out the Accept method")
//             },
//             InitializeFunc: func(ctx context.Context, t bus.Transport) error {
// 	               panic("TODO: mock out the Initialize method")
//             },
//         }
//
//         // TODO: use mockedInboundPipeline in code that requires InboundPipeline
//         //       and then make assertions.
//
//     }
type InboundPipelineMock struct {
	// AcceptFunc mocks the Accept method.
	AcceptFunc func(in1 context.Context, in2 bus.MessageSink, in3 bus.InboundEnvelope) error

	// InitializeFunc mocks the Initialize method.
	InitializeFunc func(ctx context.Context, t bus.Transport) error

	// calls tracks calls to the methods.
	calls struct {
		// Accept holds details about calls to the Accept method.
		Accept []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 bus.MessageSink
			// In3 is the in3 argument value.
			In3 bus.InboundEnvelope
		}
		// Initialize holds details about calls to the Initialize method.
		Initialize []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// T is the t argument value.
			T bus.Transport
		}
	}
}

// Accept calls AcceptFunc.
func (mock *InboundPipelineMock) Accept(in1 context.Context, in2 bus.MessageSink, in3 bus.InboundEnvelope) error {
	if mock.AcceptFunc == nil {
		panic("moq: InboundPipelineMock.AcceptFunc is nil but InboundPipeline.Accept was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 bus.MessageSink
		In3 bus.InboundEnvelope
	}{
		In1: in1,
		In2: in2,
		In3: in3,
	}
	lockInboundPipelineMockAccept.Lock()
	mock.calls.Accept = append(mock.calls.Accept, callInfo)
	lockInboundPipelineMockAccept.Unlock()
	return mock.AcceptFunc(in1, in2, in3)
}

// AcceptCalls gets all the calls that were made to Accept.
// Check the length with:
//     len(mockedInboundPipeline.AcceptCalls())
func (mock *InboundPipelineMock) AcceptCalls() []struct {
	In1 context.Context
	In2 bus.MessageSink
	In3 bus.InboundEnvelope
} {
	var calls []struct {
		In1 context.Context
		In2 bus.MessageSink
		In3 bus.InboundEnvelope
	}
	lockInboundPipelineMockAccept.RLock()
	calls = mock.calls.Accept
	lockInboundPipelineMockAccept.RUnlock()
	return calls
}

// Initialize calls InitializeFunc.
func (mock *InboundPipelineMock) Initialize(ctx context.Context, t bus.Transport) error {
	if mock.InitializeFunc == nil {
		panic("moq: InboundPipelineMock.InitializeFunc is nil but InboundPipeline.Initialize was just called")
	}
	callInfo := struct {
		Ctx context.Context
		T   bus.Transport
	}{
		Ctx: ctx,
		T:   t,
	}
	lockInboundPipelineMockInitialize.Lock()
	mock.calls.Initialize = append(mock.calls.Initialize, callInfo)
	lockInboundPipelineMockInitialize.Unlock()
	return mock.InitializeFunc(ctx, t)
}

// InitializeCalls gets all the calls that were made to Initialize.
// Check the length with:
//     len(mockedInboundPipeline.InitializeCalls())
func (mock *InboundPipelineMock) InitializeCalls() []struct {
	Ctx context.Context
	T   bus.Transport
} {
	var calls []struct {
		Ctx context.Context
		T   bus.Transport
	}
	lockInboundPipelineMockInitialize.RLock()
	calls = mock.calls.Initialize
	lockInboundPipelineMockInitialize.RUnlock()
	return calls
}

var (
	lockOutboundPipelineMockAccept     sync.RWMutex
	lockOutboundPipelineMockInitialize sync.RWMutex
)

// OutboundPipelineMock is a mock implementation of OutboundPipeline.
//
//     func TestSomethingThatUsesOutboundPipeline(t *testing.T) {
//
//         // make and configure a mocked OutboundPipeline
//         mockedOutboundPipeline := &OutboundPipelineMock{
//             AcceptFunc: func(ctx context.Context, env bus.OutboundEnvelope) error {
// 	               panic("TODO: mock out the Accept method")
//             },
//             InitializeFunc: func(ctx context.Context, t bus.Transport) error {
// 	               panic("TODO: mock out the Initialize method")
//             },
//         }
//
//         // TODO: use mockedOutboundPipeline in code that requires OutboundPipeline
//         //       and then make assertions.
//
//     }
type OutboundPipelineMock struct {
	// AcceptFunc mocks the Accept method.
	AcceptFunc func(ctx context.Context, env bus.OutboundEnvelope) error

	// InitializeFunc mocks the Initialize method.
	InitializeFunc func(ctx context.Context, t bus.Transport) error

	// calls tracks calls to the methods.
	calls struct {
		// Accept holds details about calls to the Accept method.
		Accept []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Env is the env argument value.
			Env bus.OutboundEnvelope
		}
		// Initialize holds details about calls to the Initialize method.
		Initialize []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// T is the t argument value.
			T bus.Transport
		}
	}
}

// Accept calls AcceptFunc.
func (mock *OutboundPipelineMock) Accept(ctx context.Context, env bus.OutboundEnvelope) error {
	if mock.AcceptFunc == nil {
		panic("moq: OutboundPipelineMock.AcceptFunc is nil but OutboundPipeline.Accept was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Env bus.OutboundEnvelope
	}{
		Ctx: ctx,
		Env: env,
	}
	lockOutboundPipelineMockAccept.Lock()
	mock.calls.Accept = append(mock.calls.Accept, callInfo)
	lockOutboundPipelineMockAccept.Unlock()
	return mock.AcceptFunc(ctx, env)
}

// AcceptCalls gets all the calls that were made to Accept.
// Check the length with:
//     len(mockedOutboundPipeline.AcceptCalls())
func (mock *OutboundPipelineMock) AcceptCalls() []struct {
	Ctx context.Context
	Env bus.OutboundEnvelope
} {
	var calls []struct {
		Ctx context.Context
		Env bus.OutboundEnvelope
	}
	lockOutboundPipelineMockAccept.RLock()
	calls = mock.calls.Accept
	lockOutboundPipelineMockAccept.RUnlock()
	return calls
}

// Initialize calls InitializeFunc.
func (mock *OutboundPipelineMock) Initialize(ctx context.Context, t bus.Transport) error {
	if mock.InitializeFunc == nil {
		panic("moq: OutboundPipelineMock.InitializeFunc is nil but OutboundPipeline.Initialize was just called")
	}
	callInfo := struct {
		Ctx context.Context
		T   bus.Transport
	}{
		Ctx: ctx,
		T:   t,
	}
	lockOutboundPipelineMockInitialize.Lock()
	mock.calls.Initialize = append(mock.calls.Initialize, callInfo)
	lockOutboundPipelineMockInitialize.Unlock()
	return mock.InitializeFunc(ctx, t)
}

// InitializeCalls gets all the calls that were made to Initialize.
// Check the length with:
//     len(mockedOutboundPipeline.InitializeCalls())
func (mock *OutboundPipelineMock) InitializeCalls() []struct {
	Ctx context.Context
	T   bus.Transport
} {
	var calls []struct {
		Ctx context.Context
		T   bus.Transport
	}
	lockOutboundPipelineMockInitialize.RLock()
	calls = mock.calls.Initialize
	lockOutboundPipelineMockInitialize.RUnlock()
	return calls
}
