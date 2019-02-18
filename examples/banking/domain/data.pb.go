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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

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
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xbd, 0x4e, 0xc3, 0x30,
	0x10, 0x80, 0x95, 0xd2, 0xdf, 0xab, 0x4a, 0x91, 0x19, 0xc8, 0x82, 0xa8, 0x32, 0x54, 0x59, 0x88,
	0x07, 0x9e, 0xa0, 0x74, 0xea, 0x84, 0x54, 0x31, 0xb1, 0x58, 0x17, 0xc7, 0x2d, 0x86, 0xf8, 0x1c,
	0xc5, 0xae, 0xd4, 0x8d, 0x57, 0xe3, 0xd1, 0x50, 0x12, 0x97, 0x96, 0x81, 0xed, 0xf4, 0xf9, 0xb3,
	0xf4, 0xd9, 0x07, 0x89, 0x3a, 0xa2, 0xa9, 0x4a, 0xe5, 0x78, 0x8e, 0xf4, 0xa9, 0x69, 0xcf, 0x0b,
	0x6b, 0x50, 0x13, 0x2f, 0xd0, 0x63, 0x56, 0xd5, 0xd6, 0x5b, 0x76, 0x8b, 0xc7, 0xec, 0xa4, 0x65,
	0x41, 0x4b, 0xbe, 0x60, 0xb4, 0x92, 0xd2, 0x1e, 0xc8, 0xb3, 0x7b, 0x00, 0xec, 0x46, 0xa1, 0x8b,
	0x38, 0x5a, 0x44, 0xe9, 0x64, 0x3b, 0x09, 0x64, 0x53, 0xb0, 0x3b, 0x18, 0x69, 0x27, 0x6c, 0xa5,
	0x28, 0xee, 0x2d, 0xa2, 0x74, 0xbc, 0x1d, 0x6a, 0xf7, 0x52, 0x29, 0x62, 0x0c, 0xfa, 0x84, 0x46,
	0xc5, 0x57, 0xed, 0x8d, 0x76, 0x66, 0x29, 0xdc, 0xe4, 0x58, 0x22, 0x49, 0x25, 0x34, 0x09, 0xa9,
	0xc8, 0xbb, 0xb8, 0xbf, 0x88, 0xd2, 0xc1, 0xf6, 0x3a, 0xf0, 0x0d, 0xad, 0x1b, 0x9a, 0x7c, 0x47,
	0x30, 0x7e, 0xad, 0x91, 0xdc, 0x4e, 0xd5, 0xec, 0x01, 0xa6, 0x3e, 0xcc, 0xe7, 0x06, 0x38, 0xa1,
	0x4d, 0xc1, 0x96, 0x30, 0xdf, 0xd5, 0xd6, 0x88, 0x8b, 0xd0, 0x5e, 0x2b, 0xcd, 0x1a, 0xbc, 0xfa,
	0x8d, 0x4d, 0x60, 0xe6, 0xed, 0xa5, 0xd5, 0xc5, 0x4d, 0xbd, 0x3d, 0x3b, 0x4b, 0x98, 0xa3, 0xe9,
	0xce, 0xff, 0x26, 0xce, 0x3a, 0x1c, 0x0a, 0x9b, 0x28, 0xed, 0x84, 0xb4, 0xcd, 0xc7, 0x79, 0x15,
	0x0f, 0xda, 0xc7, 0x83, 0x76, 0xeb, 0x40, 0x9e, 0xf9, 0xdb, 0xe3, 0x5e, 0xfb, 0xf7, 0x43, 0x9e,
	0x49, 0x6b, 0xf8, 0x87, 0xc1, 0xb2, 0xb4, 0x92, 0xe3, 0x91, 0xff, 0xb3, 0x94, 0x7c, 0xd8, 0x2e,
	0xe4, 0xe9, 0x27, 0x00, 0x00, 0xff, 0xff, 0xd0, 0x60, 0xb3, 0xf1, 0xb6, 0x01, 0x00, 0x00,
}
