// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sagaisnapshot..proto

package saga

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import any "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// SagaSnapshot is a container for a saga snapshot serialized into the
// protobuf format before being persisted in a Bolt bucket
type SagaSnapshot struct {
	InstanceId           string   `protobuf:"bytes,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	Revision             int64    `protobuf:"varint,2,opt,name=revision,proto3" json:"revision,omitempty"`
	PersistenceKey       string   `protobuf:"bytes,3,opt,name=persistence_key,json=persistenceKey,proto3" json:"persistence_key,omitempty"`
	Data                 *any.Any `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	InsertTime           string   `protobuf:"bytes,5,opt,name=insert_time,json=insertTime,proto3" json:"insert_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SagaSnapshot) Reset()         { *m = SagaSnapshot{} }
func (m *SagaSnapshot) String() string { return proto.CompactTextString(m) }
func (*SagaSnapshot) ProtoMessage()    {}
func (*SagaSnapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_sagaisnapshot__a40d4b4486ee8e81, []int{0}
}
func (m *SagaSnapshot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SagaSnapshot.Unmarshal(m, b)
}
func (m *SagaSnapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SagaSnapshot.Marshal(b, m, deterministic)
}
func (dst *SagaSnapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SagaSnapshot.Merge(dst, src)
}
func (m *SagaSnapshot) XXX_Size() int {
	return xxx_messageInfo_SagaSnapshot.Size(m)
}
func (m *SagaSnapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_SagaSnapshot.DiscardUnknown(m)
}

var xxx_messageInfo_SagaSnapshot proto.InternalMessageInfo

func (m *SagaSnapshot) GetInstanceId() string {
	if m != nil {
		return m.InstanceId
	}
	return ""
}

func (m *SagaSnapshot) GetRevision() int64 {
	if m != nil {
		return m.Revision
	}
	return 0
}

func (m *SagaSnapshot) GetPersistenceKey() string {
	if m != nil {
		return m.PersistenceKey
	}
	return ""
}

func (m *SagaSnapshot) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SagaSnapshot) GetInsertTime() string {
	if m != nil {
		return m.InsertTime
	}
	return ""
}

func init() {
	proto.RegisterType((*SagaSnapshot)(nil), "axbolt.saga.SagaSnapshot")
}

func init() {
	proto.RegisterFile("sagaisnapshot..proto", fileDescriptor_sagaisnapshot__a40d4b4486ee8e81)
}

var fileDescriptor_sagaisnapshot__a40d4b4486ee8e81 = []byte{
	// 222 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0xcd, 0x4a, 0x03, 0x31,
	0x14, 0x85, 0x89, 0x1d, 0x8b, 0x66, 0x44, 0x21, 0x74, 0x11, 0xbb, 0x71, 0x70, 0x63, 0x56, 0x29,
	0xe8, 0x13, 0xe8, 0x4e, 0xdc, 0x4d, 0x5d, 0xb9, 0x29, 0x77, 0x9c, 0x6b, 0xbc, 0xd8, 0x26, 0x43,
	0x72, 0x15, 0xf3, 0x74, 0xbe, 0x9a, 0x4c, 0xfa, 0x43, 0x97, 0x39, 0xe7, 0x7c, 0xe1, 0xbb, 0x72,
	0x96, 0xc0, 0x01, 0x25, 0x0f, 0x43, 0xfa, 0x0c, 0x6c, 0xed, 0x10, 0x03, 0x07, 0x55, 0xc3, 0x6f,
	0x17, 0xd6, 0x6c, 0xc7, 0x72, 0x7e, 0xed, 0x42, 0x70, 0x6b, 0x5c, 0x94, 0xaa, 0xfb, 0xfe, 0x58,
	0x80, 0xcf, 0xdb, 0xdd, 0xed, 0x9f, 0x90, 0x17, 0x4b, 0x70, 0xb0, 0xdc, 0xf1, 0xea, 0x46, 0xd6,
	0xe4, 0x13, 0x83, 0x7f, 0xc7, 0x15, 0xf5, 0x5a, 0x34, 0xc2, 0x9c, 0xb7, 0x72, 0x1f, 0x3d, 0xf7,
	0x6a, 0x2e, 0xcf, 0x22, 0xfe, 0x50, 0xa2, 0xe0, 0xf5, 0x49, 0x23, 0xcc, 0xa4, 0x3d, 0xbc, 0xd5,
	0x9d, 0xbc, 0x1a, 0x30, 0x26, 0x4a, 0x8c, 0x23, 0xff, 0x85, 0x59, 0x4f, 0xca, 0x07, 0x97, 0x47,
	0xf1, 0x0b, 0x66, 0x65, 0x64, 0xd5, 0x03, 0x83, 0xae, 0x1a, 0x61, 0xea, 0xfb, 0x99, 0xdd, 0x0a,
	0xda, 0xbd, 0xa0, 0x7d, 0xf4, 0xb9, 0x2d, 0x8b, 0x9d, 0x0f, 0x46, 0x5e, 0x31, 0x6d, 0x50, 0x9f,
	0x1e, 0x7c, 0x30, 0xf2, 0x2b, 0x6d, 0xf0, 0x69, 0xfa, 0x56, 0x8d, 0x47, 0x76, 0xd3, 0x02, 0x3f,
	0xfc, 0x07, 0x00, 0x00, 0xff, 0xff, 0x7a, 0x7c, 0xaa, 0xee, 0x10, 0x01, 0x00, 0x00,
}
