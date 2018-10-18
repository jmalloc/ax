package routing

import (
	"context"
	"reflect"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/typeswitch"
)

// MessageHandler is an interface for application-defined message handlers.
//
// Message handlers are typically the last stage in the inbound message
// pipeline. Each message handler declares its interest in a specific set
// of message types and is notified when any matching message arrives.
type MessageHandler interface {
	// MessageTypes returns the set of messages that the handler intends
	// to handle.
	//
	// The return value should be constant as it may be cached by various
	// independent stages in the message pipeline.
	MessageTypes() ax.MessageTypeSet

	// HandleMessage invokes application-defined logic that handles a
	// message.
	//
	// It may panic if env.Message is not one of the types described by
	// MessageTypes().
	HandleMessage(ctx context.Context, mctx ax.MessageContext) error
}

// NewMessageHandler returns a new message handler that dispatches messages to
// methods on an arbitrary value.
//
// For each message type to be handled, the value must implement a "handler"
// method that adheres to one of the following signatures:
//
//     func (msg *<T>)
//     func (mctx MessageContext, msg *<T>)
//     func (ctx context.Context, msg *<T>) error
//     func (ctx context.Context, mctx MessageContext, msg *<T>) error
//
// Where T is a struct type that implements ax.Message.
//
// The names of handler methods are not meaningful. By convention the methods
// are named the same as the message they accept, such as:
//
//     func (*BankAccount) CreditAccount(*messages.CreditAccount)
func NewMessageHandler(v interface{}) MessageHandler {
	sw, _, err := typeswitch.New(
		[]reflect.Type{
			reflect.TypeOf(v),
			reflect.TypeOf((*ax.Message)(nil)).Elem(),
			reflect.TypeOf((*context.Context)(nil)).Elem(),
			reflect.TypeOf((*ax.MessageContext)(nil)).Elem(),
		},
		[]reflect.Type{
			reflect.TypeOf((*error)(nil)).Elem(),
		},
		handlerSignature,
		handlerSignatureWithMessageContext,
		handlerSignatureWithError,
		handlerSignatureWithErrorAndMessageContext,
	)
	if err != nil {
		panic(err)
	}

	return &messageHandler{
		value:  v,
		types:  ax.TypesByGoType(sw.Types()...),
		handle: sw,
	}
}

type messageHandler struct {
	value  interface{}
	types  ax.MessageTypeSet
	handle typeswitch.Switch
}

// MessageTypes returns the set of messages that the handler intends
// to handle.
//
// The return value should be constant as it may be cached by various
// independent stages in the message pipeline.
func (h *messageHandler) MessageTypes() ax.MessageTypeSet {
	return h.types
}

// HandleMessage invokes application-defined logic that handles a
// message.
//
// It may panic if env.Message is not one of the types described by
// MessageTypes().
func (h *messageHandler) HandleMessage(ctx context.Context, mctx ax.MessageContext) error {
	out := h.handle.Dispatch(
		h.value,
		mctx.Envelope.Message,
		ctx,
		mctx,
	)

	if err := out[0]; err != nil {
		return err.(error)
	}

	return nil
}

var (
	handlerSignature = &typeswitch.Signature{
		In: []reflect.Type{
			reflect.TypeOf((*interface{})(nil)).Elem(),
			reflect.TypeOf((*ax.Message)(nil)).Elem(),
		},
	}

	handlerSignatureWithMessageContext = &typeswitch.Signature{
		In: []reflect.Type{
			reflect.TypeOf((*interface{})(nil)).Elem(),
			reflect.TypeOf((*ax.MessageContext)(nil)).Elem(),
			reflect.TypeOf((*ax.Message)(nil)).Elem(),
		},
	}

	handlerSignatureWithError = &typeswitch.Signature{
		In: []reflect.Type{
			reflect.TypeOf((*interface{})(nil)).Elem(),
			reflect.TypeOf((*context.Context)(nil)).Elem(),
			reflect.TypeOf((*ax.Message)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf((*error)(nil)).Elem(),
		},
	}

	handlerSignatureWithErrorAndMessageContext = &typeswitch.Signature{
		In: []reflect.Type{
			reflect.TypeOf((*interface{})(nil)).Elem(),
			reflect.TypeOf((*context.Context)(nil)).Elem(),
			reflect.TypeOf((*ax.MessageContext)(nil)).Elem(),
			reflect.TypeOf((*ax.Message)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf((*error)(nil)).Elem(),
		},
	}
)
