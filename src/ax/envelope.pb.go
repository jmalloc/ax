// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/ax/envelope.proto

package ax

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import any "github.com/golang/protobuf/ptypes/any"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// EnvelopeProto is a Protocol Buffers representation of an Envelope.
type EnvelopeProto struct {
	MessageId            string               `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	CausationId          string               `protobuf:"bytes,2,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	CorrelationId        string               `protobuf:"bytes,3,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	SendAt               *timestamp.Timestamp `protobuf:"bytes,5,opt,name=send_at,json=sendAt,proto3" json:"send_at,omitempty"`
	Message              *any.Any             `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *EnvelopeProto) Reset()         { *m = EnvelopeProto{} }
func (m *EnvelopeProto) String() string { return proto.CompactTextString(m) }
func (*EnvelopeProto) ProtoMessage()    {}
func (*EnvelopeProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_envelope_88b6aca8b0532ac9, []int{0}
}
func (m *EnvelopeProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnvelopeProto.Unmarshal(m, b)
}
func (m *EnvelopeProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnvelopeProto.Marshal(b, m, deterministic)
}
func (dst *EnvelopeProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnvelopeProto.Merge(dst, src)
}
func (m *EnvelopeProto) XXX_Size() int {
	return xxx_messageInfo_EnvelopeProto.Size(m)
}
func (m *EnvelopeProto) XXX_DiscardUnknown() {
	xxx_messageInfo_EnvelopeProto.DiscardUnknown(m)
}

var xxx_messageInfo_EnvelopeProto proto.InternalMessageInfo

func (m *EnvelopeProto) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *EnvelopeProto) GetCausationId() string {
	if m != nil {
		return m.CausationId
	}
	return ""
}

func (m *EnvelopeProto) GetCorrelationId() string {
	if m != nil {
		return m.CorrelationId
	}
	return ""
}

func (m *EnvelopeProto) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *EnvelopeProto) GetSendAt() *timestamp.Timestamp {
	if m != nil {
		return m.SendAt
	}
	return nil
}

func (m *EnvelopeProto) GetMessage() *any.Any {
	if m != nil {
		return m.Message
	}
	return nil
}

func init() {
	proto.RegisterType((*EnvelopeProto)(nil), "ax.EnvelopeProto")
}

func init() { proto.RegisterFile("src/ax/envelope.proto", fileDescriptor_envelope_88b6aca8b0532ac9) }

var fileDescriptor_envelope_88b6aca8b0532ac9 = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x8f, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x49, 0xac, 0x29, 0x99, 0x5a, 0x0f, 0x8b, 0x42, 0x0c, 0x88, 0x55, 0x10, 0x7a, 0xda,
	0x80, 0x3d, 0x79, 0x8c, 0xe0, 0x21, 0x37, 0x09, 0x9e, 0xbc, 0x94, 0x69, 0x76, 0x0c, 0x81, 0x64,
	0x37, 0xec, 0x6e, 0x25, 0xfd, 0x26, 0x7e, 0x5c, 0xe9, 0xe6, 0x0f, 0xa2, 0x87, 0x5e, 0xdf, 0xfb,
	0xfd, 0x86, 0x37, 0x70, 0x6d, 0x74, 0x91, 0x60, 0x97, 0x90, 0xfc, 0xa2, 0x5a, 0xb5, 0xc4, 0x5b,
	0xad, 0xac, 0x62, 0x3e, 0x76, 0xf1, 0x4d, 0xa9, 0x54, 0x59, 0x53, 0xe2, 0x92, 0xdd, 0xfe, 0x33,
	0x41, 0x79, 0xe8, 0xeb, 0xf8, 0xee, 0x6f, 0x65, 0xab, 0x86, 0x8c, 0xc5, 0xa6, 0xed, 0x81, 0x87,
	0x6f, 0x1f, 0x96, 0xaf, 0xc3, 0xc9, 0x37, 0x77, 0xf1, 0x16, 0xa0, 0x21, 0x63, 0xb0, 0xa4, 0x6d,
	0x25, 0x22, 0x6f, 0xe5, 0xad, 0xc3, 0x3c, 0x1c, 0x92, 0x4c, 0xb0, 0x7b, 0xb8, 0x28, 0x70, 0x6f,
	0xd0, 0x56, 0x4a, 0x1e, 0x01, 0xdf, 0x01, 0x8b, 0x29, 0xcb, 0x04, 0x7b, 0x84, 0xcb, 0x42, 0x69,
	0x4d, 0xf5, 0x04, 0x9d, 0x39, 0x68, 0xf9, 0x2b, 0xcd, 0x04, 0x7b, 0x06, 0x28, 0x34, 0xa1, 0x25,
	0xb1, 0x45, 0x1b, 0xcd, 0x56, 0xde, 0x7a, 0xf1, 0x14, 0xf3, 0x7e, 0x30, 0x1f, 0x07, 0xf3, 0xf7,
	0x71, 0x70, 0x1e, 0x0e, 0x74, 0x6a, 0xd9, 0x06, 0xe6, 0x86, 0xa4, 0xf3, 0xce, 0x4f, 0x7a, 0xc1,
	0x11, 0x4d, 0x2d, 0xe3, 0x30, 0x1f, 0xde, 0x88, 0x02, 0x27, 0x5d, 0xfd, 0x93, 0x52, 0x79, 0xc8,
	0x47, 0xe8, 0x65, 0xf6, 0xe1, 0x63, 0xb7, 0x0b, 0x5c, 0xb9, 0xf9, 0x09, 0x00, 0x00, 0xff, 0xff,
	0x6e, 0xe8, 0xdf, 0x84, 0x80, 0x01, 0x00, 0x00,
}
