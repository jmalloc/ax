// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/jmalloc/ax/endpoint"
	"github.com/jmalloc/ax/observability"
	"sync"
)

var (
	lockObserverMockAfterInbound   sync.RWMutex
	lockObserverMockAfterOutbound  sync.RWMutex
	lockObserverMockBeforeInbound  sync.RWMutex
	lockObserverMockBeforeOutbound sync.RWMutex
)

// Ensure, that ObserverMock does implement Observer.
// If this is not the case, regenerate this file with moq.
var _ observability.Observer = &ObserverMock{}

// ObserverMock is a mock implementation of Observer.
//
//     func TestSomethingThatUsesObserver(t *testing.T) {
//
//         // make and configure a mocked Observer
//         mockedObserver := &ObserverMock{
//             AfterInboundFunc: func(ctx context.Context, env endpoint.InboundEnvelope, err error)  {
// 	               panic("mock out the AfterInbound method")
//             },
//             AfterOutboundFunc: func(ctx context.Context, env endpoint.OutboundEnvelope, err error)  {
// 	               panic("mock out the AfterOutbound method")
//             },
//             BeforeInboundFunc: func(ctx context.Context, env endpoint.InboundEnvelope)  {
// 	               panic("mock out the BeforeInbound method")
//             },
//             BeforeOutboundFunc: func(ctx context.Context, env endpoint.OutboundEnvelope)  {
// 	               panic("mock out the BeforeOutbound method")
//             },
//         }
//
//         // use mockedObserver in code that requires Observer
//         // and then make assertions.
//
//     }
type ObserverMock struct {
	// AfterInboundFunc mocks the AfterInbound method.
	AfterInboundFunc func(ctx context.Context, env endpoint.InboundEnvelope, err error)

	// AfterOutboundFunc mocks the AfterOutbound method.
	AfterOutboundFunc func(ctx context.Context, env endpoint.OutboundEnvelope, err error)

	// BeforeInboundFunc mocks the BeforeInbound method.
	BeforeInboundFunc func(ctx context.Context, env endpoint.InboundEnvelope)

	// BeforeOutboundFunc mocks the BeforeOutbound method.
	BeforeOutboundFunc func(ctx context.Context, env endpoint.OutboundEnvelope)

	// calls tracks calls to the methods.
	calls struct {
		// AfterInbound holds details about calls to the AfterInbound method.
		AfterInbound []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Env is the env argument value.
			Env endpoint.InboundEnvelope
			// Err is the err argument value.
			Err error
		}
		// AfterOutbound holds details about calls to the AfterOutbound method.
		AfterOutbound []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Env is the env argument value.
			Env endpoint.OutboundEnvelope
			// Err is the err argument value.
			Err error
		}
		// BeforeInbound holds details about calls to the BeforeInbound method.
		BeforeInbound []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Env is the env argument value.
			Env endpoint.InboundEnvelope
		}
		// BeforeOutbound holds details about calls to the BeforeOutbound method.
		BeforeOutbound []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Env is the env argument value.
			Env endpoint.OutboundEnvelope
		}
	}
}

// AfterInbound calls AfterInboundFunc.
func (mock *ObserverMock) AfterInbound(ctx context.Context, env endpoint.InboundEnvelope, err error) {
	if mock.AfterInboundFunc == nil {
		panic("ObserverMock.AfterInboundFunc: method is nil but Observer.AfterInbound was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Env endpoint.InboundEnvelope
		Err error
	}{
		Ctx: ctx,
		Env: env,
		Err: err,
	}
	lockObserverMockAfterInbound.Lock()
	mock.calls.AfterInbound = append(mock.calls.AfterInbound, callInfo)
	lockObserverMockAfterInbound.Unlock()
	mock.AfterInboundFunc(ctx, env, err)
}

// AfterInboundCalls gets all the calls that were made to AfterInbound.
// Check the length with:
//     len(mockedObserver.AfterInboundCalls())
func (mock *ObserverMock) AfterInboundCalls() []struct {
	Ctx context.Context
	Env endpoint.InboundEnvelope
	Err error
} {
	var calls []struct {
		Ctx context.Context
		Env endpoint.InboundEnvelope
		Err error
	}
	lockObserverMockAfterInbound.RLock()
	calls = mock.calls.AfterInbound
	lockObserverMockAfterInbound.RUnlock()
	return calls
}

// AfterOutbound calls AfterOutboundFunc.
func (mock *ObserverMock) AfterOutbound(ctx context.Context, env endpoint.OutboundEnvelope, err error) {
	if mock.AfterOutboundFunc == nil {
		panic("ObserverMock.AfterOutboundFunc: method is nil but Observer.AfterOutbound was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
		Err error
	}{
		Ctx: ctx,
		Env: env,
		Err: err,
	}
	lockObserverMockAfterOutbound.Lock()
	mock.calls.AfterOutbound = append(mock.calls.AfterOutbound, callInfo)
	lockObserverMockAfterOutbound.Unlock()
	mock.AfterOutboundFunc(ctx, env, err)
}

// AfterOutboundCalls gets all the calls that were made to AfterOutbound.
// Check the length with:
//     len(mockedObserver.AfterOutboundCalls())
func (mock *ObserverMock) AfterOutboundCalls() []struct {
	Ctx context.Context
	Env endpoint.OutboundEnvelope
	Err error
} {
	var calls []struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
		Err error
	}
	lockObserverMockAfterOutbound.RLock()
	calls = mock.calls.AfterOutbound
	lockObserverMockAfterOutbound.RUnlock()
	return calls
}

// BeforeInbound calls BeforeInboundFunc.
func (mock *ObserverMock) BeforeInbound(ctx context.Context, env endpoint.InboundEnvelope) {
	if mock.BeforeInboundFunc == nil {
		panic("ObserverMock.BeforeInboundFunc: method is nil but Observer.BeforeInbound was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Env endpoint.InboundEnvelope
	}{
		Ctx: ctx,
		Env: env,
	}
	lockObserverMockBeforeInbound.Lock()
	mock.calls.BeforeInbound = append(mock.calls.BeforeInbound, callInfo)
	lockObserverMockBeforeInbound.Unlock()
	mock.BeforeInboundFunc(ctx, env)
}

// BeforeInboundCalls gets all the calls that were made to BeforeInbound.
// Check the length with:
//     len(mockedObserver.BeforeInboundCalls())
func (mock *ObserverMock) BeforeInboundCalls() []struct {
	Ctx context.Context
	Env endpoint.InboundEnvelope
} {
	var calls []struct {
		Ctx context.Context
		Env endpoint.InboundEnvelope
	}
	lockObserverMockBeforeInbound.RLock()
	calls = mock.calls.BeforeInbound
	lockObserverMockBeforeInbound.RUnlock()
	return calls
}

// BeforeOutbound calls BeforeOutboundFunc.
func (mock *ObserverMock) BeforeOutbound(ctx context.Context, env endpoint.OutboundEnvelope) {
	if mock.BeforeOutboundFunc == nil {
		panic("ObserverMock.BeforeOutboundFunc: method is nil but Observer.BeforeOutbound was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
	}{
		Ctx: ctx,
		Env: env,
	}
	lockObserverMockBeforeOutbound.Lock()
	mock.calls.BeforeOutbound = append(mock.calls.BeforeOutbound, callInfo)
	lockObserverMockBeforeOutbound.Unlock()
	mock.BeforeOutboundFunc(ctx, env)
}

// BeforeOutboundCalls gets all the calls that were made to BeforeOutbound.
// Check the length with:
//     len(mockedObserver.BeforeOutboundCalls())
func (mock *ObserverMock) BeforeOutboundCalls() []struct {
	Ctx context.Context
	Env endpoint.OutboundEnvelope
} {
	var calls []struct {
		Ctx context.Context
		Env endpoint.OutboundEnvelope
	}
	lockObserverMockBeforeOutbound.RLock()
	calls = mock.calls.BeforeOutbound
	lockObserverMockBeforeOutbound.RUnlock()
	return calls
}