// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sagainstance.proto

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

// SagaInstance is a container for a saga instance serialized into the
// protobuf format before being persisted in a Bolt bucket
type SagaInstance struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Revision             int64    `protobuf:"varint,2,opt,name=revision,proto3" json:"revision,omitempty"`
	PersistenceKey       string   `protobuf:"bytes,3,opt,name=persistence_key,json=persistenceKey,proto3" json:"persistence_key,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Data                 *any.Any `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	InsertTime           string   `protobuf:"bytes,6,opt,name=insert_time,json=insertTime,proto3" json:"insert_time,omitempty"`
	UpdateTime           string   `protobuf:"bytes,7,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SagaInstance) Reset()         { *m = SagaInstance{} }
func (m *SagaInstance) String() string { return proto.CompactTextString(m) }
func (*SagaInstance) ProtoMessage()    {}
func (*SagaInstance) Descriptor() ([]byte, []int) {
	return fileDescriptor_sagainstance_32dcc828a59a0cba, []int{0}
}
func (m *SagaInstance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SagaInstance.Unmarshal(m, b)
}
func (m *SagaInstance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SagaInstance.Marshal(b, m, deterministic)
}
func (dst *SagaInstance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SagaInstance.Merge(dst, src)
}
func (m *SagaInstance) XXX_Size() int {
	return xxx_messageInfo_SagaInstance.Size(m)
}
func (m *SagaInstance) XXX_DiscardUnknown() {
	xxx_messageInfo_SagaInstance.DiscardUnknown(m)
}

var xxx_messageInfo_SagaInstance proto.InternalMessageInfo

func (m *SagaInstance) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SagaInstance) GetRevision() int64 {
	if m != nil {
		return m.Revision
	}
	return 0
}

func (m *SagaInstance) GetPersistenceKey() string {
	if m != nil {
		return m.PersistenceKey
	}
	return ""
}

func (m *SagaInstance) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SagaInstance) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SagaInstance) GetInsertTime() string {
	if m != nil {
		return m.InsertTime
	}
	return ""
}

func (m *SagaInstance) GetUpdateTime() string {
	if m != nil {
		return m.UpdateTime
	}
	return ""
}

func init() {
	proto.RegisterType((*SagaInstance)(nil), "axbolt.saga.SagaInstance")
}

func init() { proto.RegisterFile("sagainstance.proto", fileDescriptor_sagainstance_32dcc828a59a0cba) }

var fileDescriptor_sagainstance_32dcc828a59a0cba = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x14, 0x45, 0xe5, 0x34, 0x04, 0x70, 0x50, 0x91, 0x2c, 0x06, 0xd3, 0x85, 0x88, 0x85, 0x4c, 0xae,
	0x04, 0x5f, 0x00, 0x1b, 0x62, 0x0b, 0x4c, 0x2c, 0xd5, 0x4b, 0xfc, 0x88, 0x9e, 0x68, 0xed, 0xc8,
	0x76, 0x11, 0xf9, 0x6c, 0xfe, 0x00, 0xd9, 0x86, 0xaa, 0xa3, 0xcf, 0x3b, 0xf7, 0xea, 0xca, 0x5c,
	0x78, 0x18, 0x81, 0x8c, 0x0f, 0x60, 0x06, 0x54, 0x93, 0xb3, 0xc1, 0x8a, 0x1a, 0xbe, 0x7b, 0xbb,
	0x0d, 0x2a, 0x9e, 0x56, 0xd7, 0xa3, 0xb5, 0xe3, 0x16, 0xd7, 0xe9, 0xd4, 0xef, 0x3f, 0xd6, 0x60,
	0xe6, 0xec, 0xdd, 0xfe, 0x30, 0x7e, 0xf1, 0x0a, 0x23, 0x3c, 0xff, 0xc5, 0xc5, 0x92, 0x17, 0xa4,
	0x25, 0x6b, 0x58, 0x7b, 0xde, 0x15, 0xa4, 0xc5, 0x8a, 0x9f, 0x39, 0xfc, 0x22, 0x4f, 0xd6, 0xc8,
	0xa2, 0x61, 0xed, 0xa2, 0x3b, 0xbc, 0xc5, 0x1d, 0xbf, 0x9c, 0xd0, 0x79, 0xf2, 0x01, 0xcd, 0x80,
	0x9b, 0x4f, 0x9c, 0xe5, 0x22, 0x05, 0x97, 0x47, 0xf8, 0x05, 0x67, 0xd1, 0xf0, 0x5a, 0xa3, 0x1f,
	0x1c, 0x4d, 0x21, 0xf6, 0x94, 0x49, 0x3a, 0x46, 0xa2, 0xe5, 0xa5, 0x86, 0x00, 0xf2, 0xa4, 0x61,
	0x6d, 0x7d, 0x7f, 0xa5, 0xf2, 0x62, 0xf5, 0xbf, 0x58, 0x3d, 0x9a, 0xb9, 0x4b, 0x86, 0xb8, 0xe1,
	0x35, 0x19, 0x8f, 0x2e, 0x6c, 0x02, 0xed, 0x50, 0x56, 0xa9, 0x8b, 0x67, 0xf4, 0x46, 0x3b, 0x8c,
	0xc2, 0x7e, 0xd2, 0x10, 0x30, 0x0b, 0xa7, 0x59, 0xc8, 0x28, 0x0a, 0x4f, 0xd5, 0x7b, 0x19, 0xbf,
	0xa5, 0xaf, 0x52, 0xfb, 0xc3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe0, 0xc9, 0x38, 0x41, 0x40,
	0x01, 0x00, 0x00,
}
