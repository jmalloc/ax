// Code generated by protoc-gen-go. DO NOT EDIT.
// source: examples/banking/workflows/data.proto

package workflows

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

type Transfer struct {
	TransferId           string   `protobuf:"bytes,1,opt,name=transfer_id,json=transferId" json:"transfer_id,omitempty"`
	FromAccountId        string   `protobuf:"bytes,2,opt,name=from_account_id,json=fromAccountId" json:"from_account_id,omitempty"`
	ToAccountId          string   `protobuf:"bytes,3,opt,name=to_account_id,json=toAccountId" json:"to_account_id,omitempty"`
	AmountInCents        int32    `protobuf:"varint,4,opt,name=amount_in_cents,json=amountInCents" json:"amount_in_cents,omitempty"`
	IsApproved           bool     `protobuf:"varint,5,opt,name=is_approved,json=isApproved" json:"is_approved,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transfer) Reset()         { *m = Transfer{} }
func (m *Transfer) String() string { return proto.CompactTextString(m) }
func (*Transfer) ProtoMessage()    {}
func (*Transfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_data_2809d817cb4b9095, []int{0}
}
func (m *Transfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transfer.Unmarshal(m, b)
}
func (m *Transfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transfer.Marshal(b, m, deterministic)
}
func (dst *Transfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transfer.Merge(dst, src)
}
func (m *Transfer) XXX_Size() int {
	return xxx_messageInfo_Transfer.Size(m)
}
func (m *Transfer) XXX_DiscardUnknown() {
	xxx_messageInfo_Transfer.DiscardUnknown(m)
}

var xxx_messageInfo_Transfer proto.InternalMessageInfo

func (m *Transfer) GetTransferId() string {
	if m != nil {
		return m.TransferId
	}
	return ""
}

func (m *Transfer) GetFromAccountId() string {
	if m != nil {
		return m.FromAccountId
	}
	return ""
}

func (m *Transfer) GetToAccountId() string {
	if m != nil {
		return m.ToAccountId
	}
	return ""
}

func (m *Transfer) GetAmountInCents() int32 {
	if m != nil {
		return m.AmountInCents
	}
	return 0
}

func (m *Transfer) GetIsApproved() bool {
	if m != nil {
		return m.IsApproved
	}
	return false
}

func init() {
	proto.RegisterType((*Transfer)(nil), "ax.examples.banking.workflows.Transfer")
}

func init() {
	proto.RegisterFile("examples/banking/workflows/data.proto", fileDescriptor_data_2809d817cb4b9095)
}

var fileDescriptor_data_2809d817cb4b9095 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xcf, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0x71, 0x19, 0x28, 0x6a, 0x2f, 0x8a, 0x2a, 0x79, 0xca, 0x82, 0x88, 0x2a, 0x51, 0x65,
	0x72, 0x06, 0x9e, 0xa0, 0x30, 0x65, 0x8d, 0x98, 0x58, 0xac, 0x6b, 0xec, 0x20, 0xab, 0x8d, 0xcf,
	0xb2, 0x0d, 0xed, 0xe3, 0xf1, 0x68, 0xc8, 0x49, 0x5a, 0xba, 0x59, 0x7f, 0xff, 0x86, 0xfb, 0xe0,
	0x45, 0x9f, 0x71, 0x70, 0x47, 0x1d, 0xea, 0x3d, 0xda, 0x83, 0xb1, 0x5f, 0xf5, 0x89, 0xfc, 0xa1,
	0x3f, 0xd2, 0x29, 0xd4, 0x0a, 0x23, 0x0a, 0xe7, 0x29, 0x12, 0x7f, 0xc2, 0xb3, 0xb8, 0x48, 0x31,
	0x4b, 0x71, 0x95, 0x9b, 0x5f, 0x06, 0xcb, 0x0f, 0x8f, 0x36, 0xf4, 0xda, 0xf3, 0x67, 0xc8, 0xe2,
	0xfc, 0x96, 0x46, 0x15, 0xac, 0x64, 0xd5, 0xaa, 0x85, 0x4b, 0x6a, 0x14, 0xdf, 0xc2, 0xba, 0xf7,
	0x34, 0x48, 0xec, 0x3a, 0xfa, 0xb6, 0x31, 0xa1, 0xbb, 0x11, 0xe5, 0x29, 0xef, 0xa6, 0xda, 0x28,
	0xbe, 0x81, 0x3c, 0xd2, 0xad, 0xba, 0x1f, 0x55, 0x16, 0xe9, 0xdf, 0x6c, 0x61, 0x8d, 0xc3, 0xf4,
	0x6f, 0x65, 0xa7, 0x6d, 0x0c, 0xc5, 0x43, 0xc9, 0xaa, 0x45, 0x9b, 0x4f, 0xb9, 0xb1, 0xef, 0x29,
	0xa6, 0xa3, 0x4c, 0x90, 0xe8, 0x9c, 0xa7, 0x1f, 0xad, 0x8a, 0x45, 0xc9, 0xaa, 0x65, 0x0b, 0x26,
	0xec, 0xe6, 0xf2, 0x96, 0x7d, 0xae, 0xae, 0x7b, 0xf6, 0x8f, 0xe3, 0xea, 0xd7, 0xbf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xf9, 0xab, 0x27, 0x4c, 0x1e, 0x01, 0x00, 0x00,
}
