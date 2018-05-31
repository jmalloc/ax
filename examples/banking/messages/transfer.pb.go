// Code generated by protoc-gen-go. DO NOT EDIT.
// source: examples/banking/messages/transfer.proto

package messages

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

// StartTransfer is a command that starts a new funds transfer.
type StartTransfer struct {
	TransferId           string   `protobuf:"bytes,1,opt,name=transfer_id,json=transferId" json:"transfer_id,omitempty"`
	FromAccountId        string   `protobuf:"bytes,2,opt,name=from_account_id,json=fromAccountId" json:"from_account_id,omitempty"`
	ToAccountId          string   `protobuf:"bytes,3,opt,name=to_account_id,json=toAccountId" json:"to_account_id,omitempty"`
	AmountInCents        int32    `protobuf:"varint,4,opt,name=amount_in_cents,json=amountInCents" json:"amount_in_cents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartTransfer) Reset()         { *m = StartTransfer{} }
func (m *StartTransfer) String() string { return proto.CompactTextString(m) }
func (*StartTransfer) ProtoMessage()    {}
func (*StartTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_transfer_8743ef8a4f199f98, []int{0}
}
func (m *StartTransfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartTransfer.Unmarshal(m, b)
}
func (m *StartTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartTransfer.Marshal(b, m, deterministic)
}
func (dst *StartTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartTransfer.Merge(dst, src)
}
func (m *StartTransfer) XXX_Size() int {
	return xxx_messageInfo_StartTransfer.Size(m)
}
func (m *StartTransfer) XXX_DiscardUnknown() {
	xxx_messageInfo_StartTransfer.DiscardUnknown(m)
}

var xxx_messageInfo_StartTransfer proto.InternalMessageInfo

func (m *StartTransfer) GetTransferId() string {
	if m != nil {
		return m.TransferId
	}
	return ""
}

func (m *StartTransfer) GetFromAccountId() string {
	if m != nil {
		return m.FromAccountId
	}
	return ""
}

func (m *StartTransfer) GetToAccountId() string {
	if m != nil {
		return m.ToAccountId
	}
	return ""
}

func (m *StartTransfer) GetAmountInCents() int32 {
	if m != nil {
		return m.AmountInCents
	}
	return 0
}

// TransferStarted is an event that occurs when a funds transfer is started.
type TransferStarted struct {
	TransferId           string   `protobuf:"bytes,1,opt,name=transfer_id,json=transferId" json:"transfer_id,omitempty"`
	FromAccountId        string   `protobuf:"bytes,2,opt,name=from_account_id,json=fromAccountId" json:"from_account_id,omitempty"`
	ToAccountId          string   `protobuf:"bytes,3,opt,name=to_account_id,json=toAccountId" json:"to_account_id,omitempty"`
	AmountInCents        int32    `protobuf:"varint,4,opt,name=amount_in_cents,json=amountInCents" json:"amount_in_cents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransferStarted) Reset()         { *m = TransferStarted{} }
func (m *TransferStarted) String() string { return proto.CompactTextString(m) }
func (*TransferStarted) ProtoMessage()    {}
func (*TransferStarted) Descriptor() ([]byte, []int) {
	return fileDescriptor_transfer_8743ef8a4f199f98, []int{1}
}
func (m *TransferStarted) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransferStarted.Unmarshal(m, b)
}
func (m *TransferStarted) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransferStarted.Marshal(b, m, deterministic)
}
func (dst *TransferStarted) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferStarted.Merge(dst, src)
}
func (m *TransferStarted) XXX_Size() int {
	return xxx_messageInfo_TransferStarted.Size(m)
}
func (m *TransferStarted) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferStarted.DiscardUnknown(m)
}

var xxx_messageInfo_TransferStarted proto.InternalMessageInfo

func (m *TransferStarted) GetTransferId() string {
	if m != nil {
		return m.TransferId
	}
	return ""
}

func (m *TransferStarted) GetFromAccountId() string {
	if m != nil {
		return m.FromAccountId
	}
	return ""
}

func (m *TransferStarted) GetToAccountId() string {
	if m != nil {
		return m.ToAccountId
	}
	return ""
}

func (m *TransferStarted) GetAmountInCents() int32 {
	if m != nil {
		return m.AmountInCents
	}
	return 0
}

// MarkTransferApproved is a command that marks a transfer as approved.
type MarkTransferApproved struct {
	TransferId           string   `protobuf:"bytes,1,opt,name=transfer_id,json=transferId" json:"transfer_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MarkTransferApproved) Reset()         { *m = MarkTransferApproved{} }
func (m *MarkTransferApproved) String() string { return proto.CompactTextString(m) }
func (*MarkTransferApproved) ProtoMessage()    {}
func (*MarkTransferApproved) Descriptor() ([]byte, []int) {
	return fileDescriptor_transfer_8743ef8a4f199f98, []int{2}
}
func (m *MarkTransferApproved) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MarkTransferApproved.Unmarshal(m, b)
}
func (m *MarkTransferApproved) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MarkTransferApproved.Marshal(b, m, deterministic)
}
func (dst *MarkTransferApproved) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarkTransferApproved.Merge(dst, src)
}
func (m *MarkTransferApproved) XXX_Size() int {
	return xxx_messageInfo_MarkTransferApproved.Size(m)
}
func (m *MarkTransferApproved) XXX_DiscardUnknown() {
	xxx_messageInfo_MarkTransferApproved.DiscardUnknown(m)
}

var xxx_messageInfo_MarkTransferApproved proto.InternalMessageInfo

func (m *MarkTransferApproved) GetTransferId() string {
	if m != nil {
		return m.TransferId
	}
	return ""
}

// TransferApproved is an event that occurs when a funds transfer is completed.
type TransferApproved struct {
	TransferId           string   `protobuf:"bytes,1,opt,name=transfer_id,json=transferId" json:"transfer_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransferApproved) Reset()         { *m = TransferApproved{} }
func (m *TransferApproved) String() string { return proto.CompactTextString(m) }
func (*TransferApproved) ProtoMessage()    {}
func (*TransferApproved) Descriptor() ([]byte, []int) {
	return fileDescriptor_transfer_8743ef8a4f199f98, []int{3}
}
func (m *TransferApproved) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransferApproved.Unmarshal(m, b)
}
func (m *TransferApproved) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransferApproved.Marshal(b, m, deterministic)
}
func (dst *TransferApproved) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferApproved.Merge(dst, src)
}
func (m *TransferApproved) XXX_Size() int {
	return xxx_messageInfo_TransferApproved.Size(m)
}
func (m *TransferApproved) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferApproved.DiscardUnknown(m)
}

var xxx_messageInfo_TransferApproved proto.InternalMessageInfo

func (m *TransferApproved) GetTransferId() string {
	if m != nil {
		return m.TransferId
	}
	return ""
}

func init() {
	proto.RegisterType((*StartTransfer)(nil), "ax.examples.banking.StartTransfer")
	proto.RegisterType((*TransferStarted)(nil), "ax.examples.banking.TransferStarted")
	proto.RegisterType((*MarkTransferApproved)(nil), "ax.examples.banking.MarkTransferApproved")
	proto.RegisterType((*TransferApproved)(nil), "ax.examples.banking.TransferApproved")
}

func init() {
	proto.RegisterFile("examples/banking/messages/transfer.proto", fileDescriptor_transfer_8743ef8a4f199f98)
}

var fileDescriptor_transfer_8743ef8a4f199f98 = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x91, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0x47, 0x65, 0xfe, 0x09, 0xae, 0xb2, 0x82, 0x0c, 0x43, 0x37, 0x2a, 0x0f, 0x55, 0xa6, 0x64,
	0xe8, 0xc0, 0x5c, 0x98, 0x32, 0xb0, 0x14, 0x26, 0x96, 0xe8, 0x1a, 0xbb, 0x55, 0x55, 0xfc, 0x47,
	0xf6, 0x81, 0xfa, 0x81, 0x10, 0x9f, 0x13, 0xc5, 0x89, 0x05, 0x23, 0x6c, 0x5d, 0xdf, 0xbd, 0x93,
	0x9e, 0xf4, 0x83, 0x52, 0x1f, 0xd0, 0xf8, 0x37, 0x1d, 0xeb, 0x35, 0xda, 0xfd, 0xce, 0x6e, 0x6b,
	0xa3, 0x63, 0xc4, 0xad, 0x8e, 0x35, 0x05, 0xb4, 0x71, 0xa3, 0x43, 0xe5, 0x83, 0x23, 0x27, 0x6e,
	0xf0, 0x50, 0x65, 0xb9, 0x1a, 0x65, 0xf9, 0xc9, 0x80, 0x3f, 0x13, 0x06, 0x7a, 0x19, 0x65, 0x71,
	0x07, 0x93, 0xfc, 0xd8, 0xee, 0xd4, 0x94, 0xcd, 0x58, 0x79, 0xb5, 0x82, 0x8c, 0x1a, 0x25, 0xe6,
	0x50, 0x6c, 0x82, 0x33, 0x2d, 0x76, 0x9d, 0x7b, 0xb7, 0xd4, 0x4b, 0x27, 0x49, 0xe2, 0x3d, 0x5e,
	0x0e, 0xb4, 0x51, 0x42, 0x02, 0x27, 0xf7, 0xdb, 0x3a, 0x4d, 0xd6, 0x84, 0xdc, 0x8f, 0x33, 0x87,
	0x02, 0xcd, 0x70, 0xb7, 0x6d, 0xa7, 0x2d, 0xc5, 0xe9, 0xd9, 0x8c, 0x95, 0xe7, 0x2b, 0x3e, 0xe0,
	0xc6, 0x3e, 0xf6, 0x50, 0x7e, 0x31, 0x28, 0x72, 0x61, 0xca, 0xd5, 0xea, 0x38, 0x43, 0xef, 0xe1,
	0xf6, 0x09, 0xc3, 0x3e, 0xb7, 0x2e, 0xbd, 0x0f, 0xee, 0xe3, 0x0f, 0xb1, 0x72, 0x01, 0xd7, 0xff,
	0x7e, 0x7a, 0x80, 0xd7, 0xcb, 0xbc, 0xf6, 0xfa, 0x22, 0xad, 0xbc, 0xf8, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0x53, 0xac, 0xe5, 0x0f, 0x11, 0x02, 0x00, 0x00,
}