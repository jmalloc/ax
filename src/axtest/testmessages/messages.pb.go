// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/axtest/testmessages/messages.proto

package testmessages

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Message is a protocol buffers message that implements ax.Message.
type Message struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Command is a protocol buffers message that implements ax.Command.
type Command struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{1}
}
func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (dst *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(dst, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// SelfValidatingCommand is a protocol buffers message that
// implements both ax.Command and endpoint.SelfValidatingMessage.
type SelfValidatingCommand struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SelfValidatingCommand) Reset()         { *m = SelfValidatingCommand{} }
func (m *SelfValidatingCommand) String() string { return proto.CompactTextString(m) }
func (*SelfValidatingCommand) ProtoMessage()    {}
func (*SelfValidatingCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{2}
}
func (m *SelfValidatingCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfValidatingCommand.Unmarshal(m, b)
}
func (m *SelfValidatingCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfValidatingCommand.Marshal(b, m, deterministic)
}
func (dst *SelfValidatingCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfValidatingCommand.Merge(dst, src)
}
func (m *SelfValidatingCommand) XXX_Size() int {
	return xxx_messageInfo_SelfValidatingCommand.Size(m)
}
func (m *SelfValidatingCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfValidatingCommand.DiscardUnknown(m)
}

var xxx_messageInfo_SelfValidatingCommand proto.InternalMessageInfo

func (m *SelfValidatingCommand) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// FailedSelfValidatingCommand is a protocol buffers message that
// implements both ax.Command and endpoint.SelfValidatingMessage.
// Its Validate method returns validation error to test failure
// scenarios in unit tests
type FailedSelfValidatingCommand struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FailedSelfValidatingCommand) Reset()         { *m = FailedSelfValidatingCommand{} }
func (m *FailedSelfValidatingCommand) String() string { return proto.CompactTextString(m) }
func (*FailedSelfValidatingCommand) ProtoMessage()    {}
func (*FailedSelfValidatingCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{3}
}
func (m *FailedSelfValidatingCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FailedSelfValidatingCommand.Unmarshal(m, b)
}
func (m *FailedSelfValidatingCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FailedSelfValidatingCommand.Marshal(b, m, deterministic)
}
func (dst *FailedSelfValidatingCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FailedSelfValidatingCommand.Merge(dst, src)
}
func (m *FailedSelfValidatingCommand) XXX_Size() int {
	return xxx_messageInfo_FailedSelfValidatingCommand.Size(m)
}
func (m *FailedSelfValidatingCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_FailedSelfValidatingCommand.DiscardUnknown(m)
}

var xxx_messageInfo_FailedSelfValidatingCommand proto.InternalMessageInfo

func (m *FailedSelfValidatingCommand) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Event is a protocol buffers message that implements ax.Event.
type Event struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{4}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// SelfValidatingEvent is a protocol buffers message that
// implements both ax.Event and endpoint.SelfValidatingMessage.
type SelfValidatingEvent struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SelfValidatingEvent) Reset()         { *m = SelfValidatingEvent{} }
func (m *SelfValidatingEvent) String() string { return proto.CompactTextString(m) }
func (*SelfValidatingEvent) ProtoMessage()    {}
func (*SelfValidatingEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{5}
}
func (m *SelfValidatingEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfValidatingEvent.Unmarshal(m, b)
}
func (m *SelfValidatingEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfValidatingEvent.Marshal(b, m, deterministic)
}
func (dst *SelfValidatingEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfValidatingEvent.Merge(dst, src)
}
func (m *SelfValidatingEvent) XXX_Size() int {
	return xxx_messageInfo_SelfValidatingEvent.Size(m)
}
func (m *SelfValidatingEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfValidatingEvent.DiscardUnknown(m)
}

var xxx_messageInfo_SelfValidatingEvent proto.InternalMessageInfo

func (m *SelfValidatingEvent) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// FailedSelfValidatingEvent is a protocol buffers message that
// implements both ax.Event and endpoint.SelfValidatingMessage.
// Its Validate method returns validation error to test failure
// scenarios in unit tests
type FailedSelfValidatingEvent struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FailedSelfValidatingEvent) Reset()         { *m = FailedSelfValidatingEvent{} }
func (m *FailedSelfValidatingEvent) String() string { return proto.CompactTextString(m) }
func (*FailedSelfValidatingEvent) ProtoMessage()    {}
func (*FailedSelfValidatingEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{6}
}
func (m *FailedSelfValidatingEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FailedSelfValidatingEvent.Unmarshal(m, b)
}
func (m *FailedSelfValidatingEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FailedSelfValidatingEvent.Marshal(b, m, deterministic)
}
func (dst *FailedSelfValidatingEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FailedSelfValidatingEvent.Merge(dst, src)
}
func (m *FailedSelfValidatingEvent) XXX_Size() int {
	return xxx_messageInfo_FailedSelfValidatingEvent.Size(m)
}
func (m *FailedSelfValidatingEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_FailedSelfValidatingEvent.DiscardUnknown(m)
}

var xxx_messageInfo_FailedSelfValidatingEvent proto.InternalMessageInfo

func (m *FailedSelfValidatingEvent) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// NonAxMessage is a protocol buffers message that does not implement ax.Message.
type NonAxMessage struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NonAxMessage) Reset()         { *m = NonAxMessage{} }
func (m *NonAxMessage) String() string { return proto.CompactTextString(m) }
func (*NonAxMessage) ProtoMessage()    {}
func (*NonAxMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_c28df44af7596237, []int{7}
}
func (m *NonAxMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NonAxMessage.Unmarshal(m, b)
}
func (m *NonAxMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NonAxMessage.Marshal(b, m, deterministic)
}
func (dst *NonAxMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NonAxMessage.Merge(dst, src)
}
func (m *NonAxMessage) XXX_Size() int {
	return xxx_messageInfo_NonAxMessage.Size(m)
}
func (m *NonAxMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_NonAxMessage.DiscardUnknown(m)
}

var xxx_messageInfo_NonAxMessage proto.InternalMessageInfo

func (m *NonAxMessage) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "axtest.testmessages.Message")
	proto.RegisterType((*Command)(nil), "axtest.testmessages.Command")
	proto.RegisterType((*SelfValidatingCommand)(nil), "axtest.testmessages.SelfValidatingCommand")
	proto.RegisterType((*FailedSelfValidatingCommand)(nil), "axtest.testmessages.FailedSelfValidatingCommand")
	proto.RegisterType((*Event)(nil), "axtest.testmessages.Event")
	proto.RegisterType((*SelfValidatingEvent)(nil), "axtest.testmessages.SelfValidatingEvent")
	proto.RegisterType((*FailedSelfValidatingEvent)(nil), "axtest.testmessages.FailedSelfValidatingEvent")
	proto.RegisterType((*NonAxMessage)(nil), "axtest.testmessages.NonAxMessage")
}

func init() {
	proto.RegisterFile("src/axtest/testmessages/messages.proto", fileDescriptor_messages_c28df44af7596237)
}

var fileDescriptor_messages_c28df44af7596237 = []byte{
	// 177 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2b, 0x2e, 0x4a, 0xd6,
	0x4f, 0xac, 0x28, 0x49, 0x2d, 0x2e, 0xd1, 0x07, 0x11, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9,
	0xc5, 0xfa, 0x30, 0x86, 0x5e, 0x41, 0x51, 0x7e, 0x49, 0xbe, 0x90, 0x30, 0x44, 0x8d, 0x1e, 0xb2,
	0x1a, 0x25, 0x79, 0x2e, 0x76, 0x5f, 0x08, 0x5b, 0x48, 0x84, 0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x34,
	0x55, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc2, 0x01, 0x29, 0x70, 0xce, 0xcf, 0xcd, 0x4d,
	0xcc, 0x4b, 0xc1, 0xa1, 0x40, 0x97, 0x4b, 0x34, 0x38, 0x35, 0x27, 0x2d, 0x2c, 0x31, 0x27, 0x33,
	0x25, 0xb1, 0x24, 0x33, 0x2f, 0x1d, 0xbf, 0x72, 0x63, 0x2e, 0x69, 0xb7, 0xc4, 0xcc, 0x9c, 0xd4,
	0x14, 0x52, 0x34, 0xc9, 0x72, 0xb1, 0xba, 0x96, 0xa5, 0xe6, 0x95, 0xe0, 0x90, 0xd6, 0xe6, 0x12,
	0x46, 0x35, 0x0d, 0x9f, 0x62, 0x43, 0x2e, 0x49, 0x6c, 0x0e, 0xc0, 0xa7, 0x45, 0x85, 0x8b, 0xc7,
	0x2f, 0x3f, 0xcf, 0xb1, 0x02, 0x6f, 0x48, 0x39, 0xf1, 0x45, 0xf1, 0x20, 0x07, 0x6d, 0x12, 0x1b,
	0x38, 0xd8, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x96, 0x73, 0x12, 0x55, 0xa0, 0x01, 0x00,
	0x00,
}
