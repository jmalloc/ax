// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: github.com/jmalloc/ax/envelope.proto

package ax

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

// EnvelopeProto is a Protocol Buffers representation of an Envelope.
type EnvelopeProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId     string                 `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	CausationId   string                 `protobuf:"bytes,2,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	CorrelationId string                 `protobuf:"bytes,3,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	SendAt        *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=send_at,json=sendAt,proto3" json:"send_at,omitempty"`
	Message       *anypb.Any             `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *EnvelopeProto) Reset() {
	*x = EnvelopeProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_jmalloc_ax_envelope_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvelopeProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvelopeProto) ProtoMessage() {}

func (x *EnvelopeProto) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_jmalloc_ax_envelope_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvelopeProto.ProtoReflect.Descriptor instead.
func (*EnvelopeProto) Descriptor() ([]byte, []int) {
	return file_github_com_jmalloc_ax_envelope_proto_rawDescGZIP(), []int{0}
}

func (x *EnvelopeProto) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *EnvelopeProto) GetCausationId() string {
	if x != nil {
		return x.CausationId
	}
	return ""
}

func (x *EnvelopeProto) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *EnvelopeProto) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *EnvelopeProto) GetSendAt() *timestamppb.Timestamp {
	if x != nil {
		return x.SendAt
	}
	return nil
}

func (x *EnvelopeProto) GetMessage() *anypb.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_github_com_jmalloc_ax_envelope_proto protoreflect.FileDescriptor

var file_github_com_jmalloc_ax_envelope_proto_rawDesc = []byte{
	0x0a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6d, 0x61,
	0x6c, 0x6c, 0x6f, 0x63, 0x2f, 0x61, 0x78, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x61, 0x78, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x98, 0x02, 0x0a, 0x0d, 0x45, 0x6e, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x75, 0x73, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63,
	0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f,
	0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x33, 0x0a, 0x07,
	0x73, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x41,
	0x74, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x42, 0x17, 0x5a, 0x15, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6a, 0x6d, 0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x2f, 0x61, 0x78, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_github_com_jmalloc_ax_envelope_proto_rawDescOnce sync.Once
	file_github_com_jmalloc_ax_envelope_proto_rawDescData = file_github_com_jmalloc_ax_envelope_proto_rawDesc
)

func file_github_com_jmalloc_ax_envelope_proto_rawDescGZIP() []byte {
	file_github_com_jmalloc_ax_envelope_proto_rawDescOnce.Do(func() {
		file_github_com_jmalloc_ax_envelope_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_jmalloc_ax_envelope_proto_rawDescData)
	})
	return file_github_com_jmalloc_ax_envelope_proto_rawDescData
}

var file_github_com_jmalloc_ax_envelope_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_jmalloc_ax_envelope_proto_goTypes = []interface{}{
	(*EnvelopeProto)(nil),         // 0: ax.EnvelopeProto
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*anypb.Any)(nil),             // 2: google.protobuf.Any
}
var file_github_com_jmalloc_ax_envelope_proto_depIdxs = []int32{
	1, // 0: ax.EnvelopeProto.created_at:type_name -> google.protobuf.Timestamp
	1, // 1: ax.EnvelopeProto.send_at:type_name -> google.protobuf.Timestamp
	2, // 2: ax.EnvelopeProto.message:type_name -> google.protobuf.Any
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_github_com_jmalloc_ax_envelope_proto_init() }
func file_github_com_jmalloc_ax_envelope_proto_init() {
	if File_github_com_jmalloc_ax_envelope_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_jmalloc_ax_envelope_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvelopeProto); i {
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
			RawDescriptor: file_github_com_jmalloc_ax_envelope_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_jmalloc_ax_envelope_proto_goTypes,
		DependencyIndexes: file_github_com_jmalloc_ax_envelope_proto_depIdxs,
		MessageInfos:      file_github_com_jmalloc_ax_envelope_proto_msgTypes,
	}.Build()
	File_github_com_jmalloc_ax_envelope_proto = out.File
	file_github_com_jmalloc_ax_envelope_proto_rawDesc = nil
	file_github_com_jmalloc_ax_envelope_proto_goTypes = nil
	file_github_com_jmalloc_ax_envelope_proto_depIdxs = nil
}
