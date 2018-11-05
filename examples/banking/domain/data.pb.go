// Code generated by protoc-gen-go. DO NOT EDIT.
// source: examples/banking/domain/data.proto

package domain

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Account contains data for the account aggregate.
type Account struct {
	AccountId            string   `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	IsOpen               bool     `protobuf:"varint,2,opt,name=is_open,json=isOpen,proto3" json:"is_open,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	BalanceInCents       int32    `protobuf:"varint,4,opt,name=balance_in_cents,json=balanceInCents,proto3" json:"balance_in_cents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_a7fadce38753b8e1, []int{0}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *Account) GetIsOpen() bool {
	if m != nil {
		return m.IsOpen
	}
	return false
}

func (m *Account) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Account) GetBalanceInCents() int32 {
	if m != nil {
		return m.BalanceInCents
	}
	return 0
}

// Transfer contains data for the transfer aggregate.
type Transfer struct {
	TransferId           string   `protobuf:"bytes,1,opt,name=transfer_id,json=transferId,proto3" json:"transfer_id,omitempty"`
	FromAccountId        string   `protobuf:"bytes,2,opt,name=from_account_id,json=fromAccountId,proto3" json:"from_account_id,omitempty"`
	ToAccountId          string   `protobuf:"bytes,3,opt,name=to_account_id,json=toAccountId,proto3" json:"to_account_id,omitempty"`
	AmountInCents        int32    `protobuf:"varint,4,opt,name=amount_in_cents,json=amountInCents,proto3" json:"amount_in_cents,omitempty"`
	IsComplete           bool     `protobuf:"varint,5,opt,name=is_complete,json=isComplete,proto3" json:"is_complete,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transfer) Reset()         { *m = Transfer{} }
func (m *Transfer) String() string { return proto.CompactTextString(m) }
func (*Transfer) ProtoMessage()    {}
func (*Transfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_a7fadce38753b8e1, []int{1}
}

func (m *Transfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transfer.Unmarshal(m, b)
}
func (m *Transfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transfer.Marshal(b, m, deterministic)
}
func (m *Transfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transfer.Merge(m, src)
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

func (m *Transfer) GetIsComplete() bool {
	if m != nil {
		return m.IsComplete
	}
	return false
}

func init() {
	proto.RegisterType((*Account)(nil), "ax.examples.banking.Account")
	proto.RegisterType((*Transfer)(nil), "ax.examples.banking.Transfer")
}

func init() { proto.RegisterFile("examples/banking/domain/data.proto", fileDescriptor_a7fadce38753b8e1) }

var fileDescriptor_a7fadce38753b8e1 = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x14, 0x45, 0x95, 0xd2, 0xa6, 0xe9, 0xab, 0x42, 0x91, 0x19, 0xc8, 0x82, 0x88, 0x32, 0x54, 0x99,
	0x92, 0x81, 0x2f, 0x28, 0x9d, 0x32, 0x21, 0x45, 0x4c, 0x2c, 0xd6, 0x4b, 0xe2, 0x22, 0x8b, 0xe6,
	0x39, 0x8a, 0x8d, 0xd4, 0x8d, 0x5f, 0xe3, 0xd3, 0x50, 0x6c, 0x97, 0x96, 0xed, 0xea, 0xf8, 0x5a,
	0x3a, 0xef, 0x42, 0x26, 0x4e, 0xd8, 0x0f, 0x47, 0xa1, 0xcb, 0x06, 0xe9, 0x53, 0xd2, 0x47, 0xd9,
	0xa9, 0x1e, 0x25, 0x95, 0x1d, 0x1a, 0x2c, 0x86, 0x51, 0x19, 0xc5, 0xee, 0xf1, 0x54, 0x9c, 0x6b,
	0x85, 0xaf, 0x65, 0xdf, 0xb0, 0xdc, 0xb5, 0xad, 0xfa, 0x22, 0xc3, 0x1e, 0x01, 0xd0, 0x45, 0x2e,
	0xbb, 0x24, 0x48, 0x83, 0x7c, 0x55, 0xaf, 0x3c, 0xa9, 0x3a, 0xf6, 0x00, 0x4b, 0xa9, 0xb9, 0x1a,
	0x04, 0x25, 0xb3, 0x34, 0xc8, 0xa3, 0x3a, 0x94, 0xfa, 0x75, 0x10, 0xc4, 0x18, 0xcc, 0x09, 0x7b,
	0x91, 0xdc, 0xd8, 0x1f, 0x36, 0xb3, 0x1c, 0xee, 0x1a, 0x3c, 0x22, 0xb5, 0x82, 0x4b, 0xe2, 0xad,
	0x20, 0xa3, 0x93, 0x79, 0x1a, 0xe4, 0x8b, 0xfa, 0xd6, 0xf3, 0x8a, 0xf6, 0x13, 0xcd, 0x7e, 0x02,
	0x88, 0xde, 0x46, 0x24, 0x7d, 0x10, 0x23, 0x7b, 0x82, 0xb5, 0xf1, 0xf9, 0xe2, 0x00, 0x67, 0x54,
	0x75, 0x6c, 0x0b, 0x9b, 0xc3, 0xa8, 0x7a, 0x7e, 0x25, 0x3a, 0xb3, 0xa5, 0x78, 0xc2, 0xbb, 0x3f,
	0xd9, 0x0c, 0x62, 0xa3, 0xae, 0x5b, 0x4e, 0x6e, 0x6d, 0xd4, 0xa5, 0xb3, 0x85, 0x0d, 0xf6, 0xee,
	0xfd, 0xbf, 0x62, 0xec, 0xb0, 0x37, 0x9c, 0xa4, 0xa4, 0xe6, 0xad, 0x9a, 0x86, 0x33, 0x22, 0x59,
	0xd8, 0xe3, 0x41, 0xea, 0xbd, 0x27, 0x2f, 0xd1, 0x7b, 0xe8, 0xd6, 0x6e, 0x42, 0xbb, 0xf4, 0xf3,
	0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x72, 0xf3, 0x64, 0x61, 0x8f, 0x01, 0x00, 0x00,
}
