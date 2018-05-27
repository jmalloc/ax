// Code generated by moq; DO NOT EDIT
// github.com/matryer/moq

package validationtest

import (
	"context"
	"sync"

	"github.com/jmalloc/ax/src/ax"
)

var (
	lockValidatorMockValidate sync.RWMutex
)

// ValidatorMock is a mock implementation of Validator.
//
//     func TestSomethingThatUsesValidator(t *testing.T) {
//
//         // make and configure a mocked Validator
//         mockedValidator := &ValidatorMock{
//             ValidateFunc: func(ctx context.Context, msg ax.Message) error {
// 	               panic("TODO: mock out the Validate method")
//             },
//         }
//
//         // TODO: use mockedValidator in code that requires Validator
//         //       and then make assertions.
//
//     }
type ValidatorMock struct {
	// ValidateFunc mocks the Validate method.
	ValidateFunc func(ctx context.Context, msg ax.Message) error

	// calls tracks calls to the methods.
	calls struct {
		// Validate holds details about calls to the Validate method.
		Validate []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Msg is the msg argument value.
			Msg ax.Message
		}
	}
}

// Validate calls ValidateFunc.
func (mock *ValidatorMock) Validate(ctx context.Context, msg ax.Message) error {
	if mock.ValidateFunc == nil {
		panic("moq: ValidatorMock.ValidateFunc is nil but Validator.Validate was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Msg ax.Message
	}{
		Ctx: ctx,
		Msg: msg,
	}
	lockValidatorMockValidate.Lock()
	mock.calls.Validate = append(mock.calls.Validate, callInfo)
	lockValidatorMockValidate.Unlock()
	return mock.ValidateFunc(ctx, msg)
}

// ValidateCalls gets all the calls that were made to Validate.
// Check the length with:
//     len(mockedValidator.ValidateCalls())
func (mock *ValidatorMock) ValidateCalls() []struct {
	Ctx context.Context
	Msg ax.Message
} {
	var calls []struct {
		Ctx context.Context
		Msg ax.Message
	}
	lockValidatorMockValidate.RLock()
	calls = mock.calls.Validate
	lockValidatorMockValidate.RUnlock()
	return calls
}
