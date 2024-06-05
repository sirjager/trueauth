// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        v3.15.8
// source: rpc_welcome.trueauth.proto

package rpc

import (
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

type WelcomeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WelcomeRequest) Reset() {
	*x = WelcomeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_welcome_trueauth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WelcomeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WelcomeRequest) ProtoMessage() {}

func (x *WelcomeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_welcome_trueauth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WelcomeRequest.ProtoReflect.Descriptor instead.
func (*WelcomeRequest) Descriptor() ([]byte, []int) {
	return file_rpc_welcome_trueauth_proto_rawDescGZIP(), []int{0}
}

type WelcomeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Docs    string `protobuf:"bytes,2,opt,name=docs,proto3" json:"docs,omitempty"`
}

func (x *WelcomeResponse) Reset() {
	*x = WelcomeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_welcome_trueauth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WelcomeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WelcomeResponse) ProtoMessage() {}

func (x *WelcomeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_welcome_trueauth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WelcomeResponse.ProtoReflect.Descriptor instead.
func (*WelcomeResponse) Descriptor() ([]byte, []int) {
	return file_rpc_welcome_trueauth_proto_rawDescGZIP(), []int{1}
}

func (x *WelcomeResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *WelcomeResponse) GetDocs() string {
	if x != nil {
		return x.Docs
	}
	return ""
}

var File_rpc_welcome_trueauth_proto protoreflect.FileDescriptor

var file_rpc_welcome_trueauth_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x72, 0x70, 0x63, 0x5f, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x2e, 0x74, 0x72,
	0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x74, 0x72,
	0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x22, 0x10, 0x0a, 0x0e, 0x57, 0x65, 0x6c, 0x63, 0x6f, 0x6d,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3f, 0x0a, 0x0f, 0x57, 0x65, 0x6c, 0x63,
	0x6f, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x6f, 0x63, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x6f, 0x63, 0x73, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x72, 0x6a, 0x61, 0x67, 0x65, 0x72,
	0x2f, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_welcome_trueauth_proto_rawDescOnce sync.Once
	file_rpc_welcome_trueauth_proto_rawDescData = file_rpc_welcome_trueauth_proto_rawDesc
)

func file_rpc_welcome_trueauth_proto_rawDescGZIP() []byte {
	file_rpc_welcome_trueauth_proto_rawDescOnce.Do(func() {
		file_rpc_welcome_trueauth_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_welcome_trueauth_proto_rawDescData)
	})
	return file_rpc_welcome_trueauth_proto_rawDescData
}

var file_rpc_welcome_trueauth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_welcome_trueauth_proto_goTypes = []interface{}{
	(*WelcomeRequest)(nil),  // 0: trueauth.WelcomeRequest
	(*WelcomeResponse)(nil), // 1: trueauth.WelcomeResponse
}
var file_rpc_welcome_trueauth_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_welcome_trueauth_proto_init() }
func file_rpc_welcome_trueauth_proto_init() {
	if File_rpc_welcome_trueauth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_welcome_trueauth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WelcomeRequest); i {
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
		file_rpc_welcome_trueauth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WelcomeResponse); i {
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
			RawDescriptor: file_rpc_welcome_trueauth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_welcome_trueauth_proto_goTypes,
		DependencyIndexes: file_rpc_welcome_trueauth_proto_depIdxs,
		MessageInfos:      file_rpc_welcome_trueauth_proto_msgTypes,
	}.Build()
	File_rpc_welcome_trueauth_proto = out.File
	file_rpc_welcome_trueauth_proto_rawDesc = nil
	file_rpc_welcome_trueauth_proto_goTypes = nil
	file_rpc_welcome_trueauth_proto_depIdxs = nil
}
