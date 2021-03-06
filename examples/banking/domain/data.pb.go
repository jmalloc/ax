// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: github.com/jmalloc/ax/examples/banking/domain/data.proto

package domain

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Account contains data for the account aggregate.
type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId      string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	IsOpen         bool   `protobuf:"varint,2,opt,name=is_open,json=isOpen,proto3" json:"is_open,omitempty"`
	Name           string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	BalanceInCents int32  `protobuf:"varint,4,opt,name=balance_in_cents,json=balanceInCents,proto3" json:"balance_in_cents,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Account) GetIsOpen() bool {
	if x != nil {
		return x.IsOpen
	}
	return false
}

func (x *Account) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Account) GetBalanceInCents() int32 {
	if x != nil {
		return x.BalanceInCents
	}
	return 0
}

// Transfer contains data for the transfer aggregate.
type Transfer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TransferId    string `protobuf:"bytes,1,opt,name=transfer_id,json=transferId,proto3" json:"transfer_id,omitempty"`
	FromAccountId string `protobuf:"bytes,2,opt,name=from_account_id,json=fromAccountId,proto3" json:"from_account_id,omitempty"`
	ToAccountId   string `protobuf:"bytes,3,opt,name=to_account_id,json=toAccountId,proto3" json:"to_account_id,omitempty"`
	AmountInCents int32  `protobuf:"varint,4,opt,name=amount_in_cents,json=amountInCents,proto3" json:"amount_in_cents,omitempty"`
	IsComplete    bool   `protobuf:"varint,5,opt,name=is_complete,json=isComplete,proto3" json:"is_complete,omitempty"`
}

func (x *Transfer) Reset() {
	*x = Transfer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transfer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transfer) ProtoMessage() {}

func (x *Transfer) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transfer.ProtoReflect.Descriptor instead.
func (*Transfer) Descriptor() ([]byte, []int) {
	return file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescGZIP(), []int{1}
}

func (x *Transfer) GetTransferId() string {
	if x != nil {
		return x.TransferId
	}
	return ""
}

func (x *Transfer) GetFromAccountId() string {
	if x != nil {
		return x.FromAccountId
	}
	return ""
}

func (x *Transfer) GetToAccountId() string {
	if x != nil {
		return x.ToAccountId
	}
	return ""
}

func (x *Transfer) GetAmountInCents() int32 {
	if x != nil {
		return x.AmountInCents
	}
	return 0
}

func (x *Transfer) GetIsComplete() bool {
	if x != nil {
		return x.IsComplete
	}
	return false
}

var File_github_com_jmalloc_ax_examples_banking_domain_data_proto protoreflect.FileDescriptor

var file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDesc = []byte{
	0x0a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6d, 0x61,
	0x6c, 0x6c, 0x6f, 0x63, 0x2f, 0x61, 0x78, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73,
	0x2f, 0x62, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x78, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x62, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x22,
	0x7f, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f,
	0x6f, 0x70, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4f, 0x70,
	0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x6e, 0x43, 0x65, 0x6e, 0x74, 0x73,
	0x22, 0xc0, 0x01, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x12, 0x1f, 0x0a,
	0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26,
	0x0a, 0x0f, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x66, 0x72, 0x6f, 0x6d, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x74, 0x6f, 0x5f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74,
	0x6f, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0d, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x6e, 0x43, 0x65, 0x6e,
	0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6a, 0x6d, 0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x2f, 0x61, 0x78, 0x2f, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x62, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescOnce sync.Once
	file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescData = file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDesc
)

func file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescGZIP() []byte {
	file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescOnce.Do(func() {
		file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescData)
	})
	return file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDescData
}

var file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_jmalloc_ax_examples_banking_domain_data_proto_goTypes = []interface{}{
	(*Account)(nil),  // 0: ax.examples.banking.Account
	(*Transfer)(nil), // 1: ax.examples.banking.Transfer
}
var file_github_com_jmalloc_ax_examples_banking_domain_data_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_jmalloc_ax_examples_banking_domain_data_proto_init() }
func file_github_com_jmalloc_ax_examples_banking_domain_data_proto_init() {
	if File_github_com_jmalloc_ax_examples_banking_domain_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transfer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_jmalloc_ax_examples_banking_domain_data_proto_goTypes,
		DependencyIndexes: file_github_com_jmalloc_ax_examples_banking_domain_data_proto_depIdxs,
		MessageInfos:      file_github_com_jmalloc_ax_examples_banking_domain_data_proto_msgTypes,
	}.Build()
	File_github_com_jmalloc_ax_examples_banking_domain_data_proto = out.File
	file_github_com_jmalloc_ax_examples_banking_domain_data_proto_rawDesc = nil
	file_github_com_jmalloc_ax_examples_banking_domain_data_proto_goTypes = nil
	file_github_com_jmalloc_ax_examples_banking_domain_data_proto_depIdxs = nil
}
