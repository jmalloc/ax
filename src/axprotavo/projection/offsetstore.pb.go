// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/axprotavo/projection/offsetstore.proto

package projection

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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_offsetstore_f49341522c654b57, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Empty)(nil), "ax.protavo.projection.Empty")
}

func init() {
	proto.RegisterFile("src/axprotavo/projection/offsetstore.proto", fileDescriptor_offsetstore_f49341522c654b57)
}

var fileDescriptor_offsetstore_f49341522c654b57 = []byte{
	// 96 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2a, 0x2e, 0x4a, 0xd6,
	0x4f, 0xac, 0x28, 0x28, 0xca, 0x2f, 0x49, 0x2c, 0xcb, 0xd7, 0x2f, 0x28, 0xca, 0xcf, 0x4a, 0x4d,
	0x2e, 0xc9, 0xcc, 0xcf, 0xd3, 0xcf, 0x4f, 0x4b, 0x2b, 0x4e, 0x2d, 0x29, 0x2e, 0xc9, 0x2f, 0x4a,
	0xd5, 0x03, 0x49, 0xe7, 0x0b, 0x89, 0x26, 0x56, 0xe8, 0x41, 0x15, 0xea, 0x21, 0x14, 0x2a, 0xb1,
	0x73, 0xb1, 0xba, 0xe6, 0x16, 0x94, 0x54, 0x3a, 0xf1, 0x44, 0x71, 0x21, 0x84, 0x93, 0xd8, 0xc0,
	0x9a, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x49, 0x6a, 0x63, 0xad, 0x62, 0x00, 0x00, 0x00,
}