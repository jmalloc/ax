package ax

import (
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
)

// MessageType provides information about a particular message type.
type MessageType struct {
	// Name is the fully-qualified Protocol Buffers name of the message type.
	Name string

	// StructType is the struct that represents messages of this type.
	//
	// Note that only a pointer-to-struct type will satisfy the Message
	// interface.
	StructType reflect.Type
}

// TypeOf returns the message type of m.
func TypeOf(m Message) MessageType {
	return MessageType{
		proto.MessageName(m),
		reflect.TypeOf(m).Elem(),
	}
}

// TypeByName returns the message type for a fully-qualified Protocol Buffers
// message name.
//
// If n is the name of a registered implementation of Message,
// then mt is the type of that message, and ok is true; otherwise, ok is false.
//
// Note that messages are only added to the registry when their respective Go
// package is imported.
func TypeByName(n string) (mt MessageType, ok bool) {
	rt := proto.MessageType(n)

	if rt == nil {
		return MessageType{}, false
	}

	if !rt.Implements(messageType) {
		return MessageType{}, false
	}

	return MessageType{
		n,
		rt.Elem(),
	}, true
}

// TypeByGoType returns the message type for the given Go type.
//
// It panics if t does not implement Message.
//
// Note that messages are only added to the registry when their respective Go
// package is imported.
func TypeByGoType(t reflect.Type) (mt MessageType) {
	v := reflect.Zero(t).Interface()

	return TypeOf(
		v.(Message),
	)
}

// ToSet returns a MessageTypeSet containing mt as its only member.
func (mt MessageType) ToSet() MessageTypeSet {
	return MessageTypeSet{
		map[MessageType]struct{}{mt: {}},
	}
}

// New returns a new pointer to a zero-value message of this type.
func (mt MessageType) New() Message {
	return reflect.New(mt.StructType).Interface().(Message)
}

// MessageName returns the Protocol Buffers message name without the package name.
func (mt MessageType) MessageName() string {
	i := strings.LastIndexByte(mt.Name, '.')
	if i == -1 {
		return mt.Name
	}

	return mt.Name[i+1:]
}

// PackageName returns the Protocol Buffers package name for this message type.
func (mt MessageType) PackageName() string {
	i := strings.LastIndexByte(mt.Name, '.')
	if i == -1 {
		return ""
	}

	return mt.Name[:i]
}

// IsCommand returns true if the message type satisfies the Command interface.
func (mt MessageType) IsCommand() bool {
	return reflect.PtrTo(mt.StructType).Implements(commandType)
}

// IsEvent returns true if the message type satisfies the Event interface.
func (mt MessageType) IsEvent() bool {
	return reflect.PtrTo(mt.StructType).Implements(eventType)
}

func (mt MessageType) String() string {
	if mt.IsCommand() {
		return mt.Name + "?"
	} else if mt.IsEvent() {
		return mt.Name + "!"
	}

	return mt.Name
}

// MessageTypeSet is a collection of unique message types.
type MessageTypeSet struct {
	members map[MessageType]struct{}
}

// NewMessageTypeSet returns a set containing the message types in mt.
func NewMessageTypeSet(mt ...MessageType) MessageTypeSet {
	members := make(map[MessageType]struct{}, len(mt))

	for _, v := range mt {
		members[v] = struct{}{}
	}

	return MessageTypeSet{members}
}

// TypesOf returns a set containing the message types of the messages in m.
func TypesOf(m ...Message) MessageTypeSet {
	members := make(map[MessageType]struct{}, len(m))

	for _, v := range m {
		members[TypeOf(v)] = struct{}{}
	}

	return MessageTypeSet{members}
}

// TypesByGoType returns a set containing the message types of the Go types in
// t.
//
// It panics if any of the types do not implement Message.
func TypesByGoType(t ...reflect.Type) MessageTypeSet {
	members := make(map[MessageType]struct{}, len(t))

	for _, v := range t {
		members[TypeByGoType(v)] = struct{}{}
	}

	return MessageTypeSet{members}
}

// Members returns the message types in the set.
func (s MessageTypeSet) Members() []MessageType {
	types := make([]MessageType, 0, len(s.members))

	for mt := range s.members {
		types = append(types, mt)
	}

	return types
}

// Len returns the number of types in the set.
func (s MessageTypeSet) Len() int {
	return len(s.members)
}

// Has returns true if mt is a member of the set.
func (s MessageTypeSet) Has(mt MessageType) bool {
	_, ok := s.members[mt]
	return ok
}

// Add returns a new set containing the members of this set and mt.
func (s MessageTypeSet) Add(mt MessageType) MessageTypeSet {
	if s.Has(mt) {
		return s
	}

	members := make(map[MessageType]struct{}, len(s.members)+1)

	members[mt] = struct{}{}

	for mt := range s.members {
		members[mt] = struct{}{}
	}

	return MessageTypeSet{members}
}

// Union returns the set union of s and o.
func (s MessageTypeSet) Union(o MessageTypeSet) MessageTypeSet {
	ol := o.Len()
	sl := s.Len()

	if ol == 0 {
		return s
	} else if sl == 0 {
		return o
	}

	members := make(map[MessageType]struct{}, ol+sl)

	for mt := range s.members {
		members[mt] = struct{}{}
	}

	for mt := range o.members {
		members[mt] = struct{}{}
	}

	return MessageTypeSet{members}
}

// Intersection returns the set intersection of s and o.
func (s MessageTypeSet) Intersection(o MessageTypeSet) MessageTypeSet {
	ol := o.Len()
	sl := s.Len()

	if ol == 0 {
		return o
	} else if sl == 0 {
		return s
	}

	// always iterate over the smaller of the two maps
	if ol < sl {
		return intersection(o, s)
	}

	return intersection(s, o)
}

func intersection(a, b MessageTypeSet) MessageTypeSet {
	members := make(map[MessageType]struct{}, len(a.members))

	for mt := range a.members {
		if _, ok := b.members[mt]; ok {
			members[mt] = struct{}{}
		}
	}

	return MessageTypeSet{members}
}

var (
	messageType = reflect.TypeOf((*Message)(nil)).Elem()
	commandType = reflect.TypeOf((*Command)(nil)).Elem()
	eventType   = reflect.TypeOf((*Event)(nil)).Elem()
)
