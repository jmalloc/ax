package aggregate

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/jmalloc/ax/src/ax"
)

// HandlerSet contains the set of handler methods for each supported message type.
type HandlerSet struct {
	// Handlers is a map of message type to the method that handles that message
	// for both commands and events
	Handlers map[ax.MessageType]reflect.Method

	// CommandTypes is the set of handlers that are available for commands.
	CommandTypes ax.MessageTypeSet
}

// NewHandlerSet returns the handler set to use for aggreates implemented by the
// given type.
func NewHandlerSet(t reflect.Type) *HandlerSet {
	hs := &HandlerSet{
		Handlers: map[ax.MessageType]reflect.Method{},
	}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		if isUnexposed(m.Name) {
			continue
		}

		if mt, ok := isHandler(m.Type); ok {
			if x, ok := hs.Handlers[mt]; ok {
				panic(fmt.Sprintf(
					"%s.%s() can be used to handle %s messages, they are already handled by %s()",
					t,
					m.Name,
					mt.Name,
					x.Name,
				))
			}

			hs.Handlers[mt] = m

			if mt.IsCommand() {
				hs.CommandTypes = hs.CommandTypes.Add(mt)
			}

		}
	}

	return hs
}

// HandleCommand invokes the command handler for the given message.
func (hs *HandlerSet) HandleCommand(agg Aggregate, m ax.Command, r Recorder) {
	mt := ax.TypeOf(m)
	hs.Handlers[mt].Func.Call(
		[]reflect.Value{
			reflect.ValueOf(agg),
			reflect.ValueOf(m),
			reflect.ValueOf(r),
		},
	)
}

// HandleEvent invokes the event handler for the given message.
func (hs *HandlerSet) HandleEvent(agg Aggregate, m ax.Event) {
	mt := ax.TypeOf(m)
	hs.Handlers[mt].Func.Call(
		[]reflect.Value{
			reflect.ValueOf(agg),
			reflect.ValueOf(m),
		},
	)
}

// isHandler returns the message type that is handled by rt, if it is a function
// that conforms to the signature for either a command or event handler.
func isHandler(rt reflect.Type) (ax.MessageType, bool) {
	if mt, ok := isCommandHandler(rt); ok {
		return mt, true
	}

	return isEventHandler(rt)
}

// isCommandHandler returns the message that is handled by rt, if it is a
// function that conforms to the signature for command handlers, which is
// func(<implementation of ax.Command>, aggregate.Recorder).
func isCommandHandler(rt reflect.Type) (ax.MessageType, bool) {
	if rt.NumIn() != 3 {
		return ax.MessageType{}, false
	}

	// arg 0 is the method receiver
	msgArg := rt.In(1)
	recArg := rt.In(2)

	if !recArg.AssignableTo(recorderType) {
		return ax.MessageType{}, false
	}

	if mt, ok := messageType(msgArg); ok {
		return mt, mt.IsCommand()
	}

	return ax.MessageType{}, false
}

// isCommandHandler returns the message that is handled by rt, if it is a
// function that conforms to the signature for command handlers, which is
// func(<implementation of ax.Event>).
func isEventHandler(rt reflect.Type) (ax.MessageType, bool) {
	if rt.NumIn() != 2 {
		return ax.MessageType{}, false
	}

	// arg 0 is the method receiver
	msgArg := rt.In(1)

	if mt, ok := messageType(msgArg); ok {
		return mt, mt.IsEvent()
	}

	return ax.MessageType{}, false
}

// messageType returns the ax.MessageType of t, if it is an ax.Message.
func messageType(t reflect.Type) (ax.MessageType, bool) {
	v := reflect.Zero(t).Interface()

	if m, ok := v.(ax.Message); ok {
		return ax.TypeOf(m), true
	}

	return ax.MessageType{}, false
}

var (
	recorderType = reflect.TypeOf(Recorder(nil))
)

// isUnexposed returns true if the method, field or type named n is unexposed.
func isUnexposed(n string) bool {
	return unicode.IsLower([]rune(n)[0])
}
