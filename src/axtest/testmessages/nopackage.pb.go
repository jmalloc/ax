// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/axtest/testmessages/nopackage.proto

package testmessages

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type NoPackage struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoPackage) Reset()         { *m = NoPackage{} }
func (m *NoPackage) String() string { return proto.CompactTextString(m) }
func (*NoPackage) ProtoMessage()    {}
func (*NoPackage) Descriptor() ([]byte, []int) {
	return fileDescriptor_6645bac7057ffc86, []int{0}
}

func (m *NoPackage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoPackage.Unmarshal(m, b)
}
func (m *NoPackage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoPackage.Marshal(b, m, deterministic)
}
func (m *NoPackage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoPackage.Merge(m, src)
}
func (m *NoPackage) XXX_Size() int {
	return xxx_messageInfo_NoPackage.Size(m)
}
func (m *NoPackage) XXX_DiscardUnknown() {
	xxx_messageInfo_NoPackage.DiscardUnknown(m)
}

var xxx_messageInfo_NoPackage proto.InternalMessageInfo

func init() {
	proto.RegisterType((*NoPackage)(nil), "NoPackage")
}

func init() {
	proto.RegisterFile("src/axtest/testmessages/nopackage.proto", fileDescriptor_6645bac7057ffc86)
}

var fileDescriptor_6645bac7057ffc86 = []byte{
	// 103 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2f, 0x2e, 0x4a, 0xd6,
	0x4f, 0xac, 0x28, 0x49, 0x2d, 0x2e, 0xd1, 0x07, 0x11, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9,
	0xc5, 0xfa, 0x79, 0xf9, 0x05, 0x89, 0xc9, 0xd9, 0x89, 0xe9, 0xa9, 0x7a, 0x05, 0x45, 0xf9, 0x25,
	0xf9, 0x4a, 0xdc, 0x5c, 0x9c, 0x7e, 0xf9, 0x01, 0x10, 0x21, 0x27, 0xfd, 0x28, 0xdd, 0xf4, 0xcc,
	0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xac, 0xdc, 0xc4, 0x9c, 0x9c, 0x7c, 0x90,
	0x31, 0xfa, 0x38, 0x4c, 0x4b, 0x62, 0x03, 0x1b, 0x62, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x97,
	0x91, 0x70, 0x24, 0x6f, 0x00, 0x00, 0x00,
}
