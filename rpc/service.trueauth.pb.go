// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        v4.25.3
// source: service.trueauth.proto

package rpc

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_trueauth_proto protoreflect.FileDescriptor

var file_service_trueauth_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75,
	0x74, 0x68, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65,
	0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x19, 0x72, 0x70, 0x63, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2e, 0x74, 0x72, 0x75,
	0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72, 0x70, 0x63,
	0x5f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x72, 0x70, 0x63, 0x5f, 0x72, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x18, 0x72, 0x70, 0x63, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x74, 0x2e, 0x74, 0x72,
	0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72, 0x70,
	0x63, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x72, 0x70, 0x63, 0x5f, 0x73, 0x69, 0x67,
	0x6e, 0x6f, 0x75, 0x74, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72, 0x70, 0x63, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x2e,
	0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17,
	0x72, 0x70, 0x63, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x72, 0x70, 0x63, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1a, 0x72, 0x70, 0x63, 0x5f, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x2e, 0x74, 0x72, 0x75,
	0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x95, 0x0c, 0x0a, 0x08,
	0x54, 0x72, 0x75, 0x65, 0x41, 0x75, 0x74, 0x68, 0x12, 0x65, 0x0a, 0x07, 0x57, 0x65, 0x6c, 0x63,
	0x6f, 0x6d, 0x65, 0x12, 0x18, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x57,
	0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x57, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x25, 0x92, 0x41, 0x19, 0x0a, 0x06, 0x53,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x0f, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x20, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x03, 0x12, 0x01, 0x2f, 0x12,
	0x66, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x17, 0x2e, 0x74, 0x72, 0x75, 0x65,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x29, 0x92, 0x41,
	0x14, 0x0a, 0x06, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x0a, 0x61, 0x70, 0x69, 0x20, 0x68,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x76, 0x31,
	0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x7b, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75,
	0x70, 0x12, 0x17, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x53, 0x69, 0x67,
	0x6e, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74, 0x72, 0x75,
	0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3e, 0x92, 0x41, 0x21, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65,
	0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0f, 0x73, 0x69, 0x67, 0x6e, 0x75,
	0x70, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x75, 0x73, 0x65, 0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14,
	0x3a, 0x01, 0x2a, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x73, 0x69,
	0x67, 0x6e, 0x75, 0x70, 0x12, 0x85, 0x01, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x12,
	0x17, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x69,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x48, 0x92, 0x41, 0x2e, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0b, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x20,
	0x75, 0x73, 0x65, 0x72, 0x62, 0x0f, 0x0a, 0x0d, 0x0a, 0x09, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41,
	0x75, 0x74, 0x68, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x12, 0x83, 0x01, 0x0a,
	0x08, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x74, 0x72, 0x75, 0x65,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x40, 0x92, 0x41, 0x21, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x73,
	0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x3a, 0x01, 0x2a, 0x22,
	0x11, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x90, 0x01, 0x0a, 0x07, 0x53, 0x69, 0x67, 0x6e, 0x6f, 0x75, 0x74, 0x12, 0x18,
	0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x6f, 0x75,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x50, 0x92, 0x41, 0x35, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e,
	0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x73, 0x69, 0x67, 0x6e, 0x6f, 0x75,
	0x74, 0x20, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x11, 0x0a, 0x0f, 0x0a, 0x0b,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x00, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x12, 0x2a, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x73, 0x69,
	0x67, 0x6e, 0x6f, 0x75, 0x74, 0x12, 0x95, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73,
	0x68, 0x12, 0x18, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x52, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x74, 0x72,
	0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x55, 0x92, 0x41, 0x3a, 0x0a, 0x0e, 0x41, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x72, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x20, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x62, 0x12, 0x0a, 0x10, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x12, 0x76, 0x0a,
	0x05, 0x52, 0x65, 0x73, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17,
	0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3c, 0x92, 0x41, 0x20, 0x0a, 0x0e, 0x41, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x20, 0x72, 0x65, 0x73, 0x65, 0x74, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x13, 0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f,
	0x72, 0x65, 0x73, 0x65, 0x74, 0x12, 0x7b, 0x0a, 0x06, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12,
	0x17, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x3e, 0x92, 0x41, 0x24, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x20, 0x76,
	0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x79, 0x12, 0x87, 0x01, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x17, 0x2e,
	0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74,
	0x68, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x4a, 0x92, 0x41, 0x32, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x20, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x69, 0x6f, 0x6e, 0x62, 0x11, 0x0a, 0x0f, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x2a, 0x0d, 0x2f,
	0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x12, 0x7c, 0x0a, 0x04,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x74, 0x72,
	0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x45, 0x92, 0x41, 0x2d, 0x0a, 0x0e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e,
	0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x08, 0x67, 0x65, 0x74, 0x20, 0x75, 0x73,
	0x65, 0x72, 0x62, 0x11, 0x0a, 0x0f, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x1a, 0x86, 0x01, 0x92, 0x41, 0x82,
	0x01, 0x12, 0x3a, 0x54, 0x72, 0x75, 0x65, 0x20, 0x41, 0x75, 0x74, 0x68, 0x20, 0x69, 0x73, 0x20,
	0x61, 0x20, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x61, 0x6c, 0x6f, 0x6e, 0x65, 0x20, 0x61, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x67, 0x52, 0x50, 0x43,
	0x20, 0x61, 0x6e, 0x64, 0x20, 0x72, 0x65, 0x73, 0x74, 0x20, 0x61, 0x70, 0x69, 0x1a, 0x44, 0x0a,
	0x1c, 0x46, 0x69, 0x6e, 0x64, 0x20, 0x6f, 0x75, 0x74, 0x20, 0x6d, 0x6f, 0x72, 0x65, 0x20, 0x61,
	0x62, 0x6f, 0x75, 0x74, 0x20, 0x54, 0x72, 0x75, 0x65, 0x41, 0x75, 0x74, 0x68, 0x12, 0x24, 0x68,
	0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x69, 0x72, 0x6a, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x72, 0x75, 0x65, 0x61,
	0x75, 0x74, 0x68, 0x42, 0xf0, 0x01, 0x92, 0x41, 0xca, 0x01, 0x12, 0x43, 0x0a, 0x08, 0x54, 0x72,
	0x75, 0x65, 0x41, 0x75, 0x74, 0x68, 0x22, 0x30, 0x0a, 0x08, 0x53, 0x69, 0x72, 0x4a, 0x61, 0x67,
	0x65, 0x72, 0x12, 0x24, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x72, 0x6a, 0x61, 0x67, 0x65, 0x72, 0x2f,
	0x74, 0x72, 0x75, 0x65, 0x61, 0x75, 0x74, 0x68, 0x32, 0x05, 0x30, 0x2e, 0x31, 0x2e, 0x30, 0x2a,
	0x03, 0x01, 0x02, 0x04, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x5a, 0x5a, 0x0a, 0x22, 0x0a, 0x0b, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x13, 0x08, 0x02, 0x1a, 0x0d, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x02, 0x0a, 0x0f,
	0x0a, 0x09, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x12, 0x02, 0x08, 0x01, 0x0a,
	0x23, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x13, 0x08, 0x02, 0x1a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x20, 0x02, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x73, 0x69, 0x72, 0x6a, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x74, 0x72, 0x75, 0x65, 0x61, 0x75,
	0x74, 0x68, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_service_trueauth_proto_goTypes = []interface{}{
	(*WelcomeRequest)(nil),   // 0: trueauth.WelcomeRequest
	(*HealthRequest)(nil),    // 1: trueauth.HealthRequest
	(*SignupRequest)(nil),    // 2: trueauth.SignupRequest
	(*SigninRequest)(nil),    // 3: trueauth.SigninRequest
	(*ValidateRequest)(nil),  // 4: trueauth.ValidateRequest
	(*SignoutRequest)(nil),   // 5: trueauth.SignoutRequest
	(*RefreshRequest)(nil),   // 6: trueauth.RefreshRequest
	(*ResetRequest)(nil),     // 7: trueauth.ResetRequest
	(*VerifyRequest)(nil),    // 8: trueauth.VerifyRequest
	(*DeleteRequest)(nil),    // 9: trueauth.DeleteRequest
	(*UserRequest)(nil),      // 10: trueauth.UserRequest
	(*WelcomeResponse)(nil),  // 11: trueauth.WelcomeResponse
	(*HealthResponse)(nil),   // 12: trueauth.HealthResponse
	(*SignupResponse)(nil),   // 13: trueauth.SignupResponse
	(*SigninResponse)(nil),   // 14: trueauth.SigninResponse
	(*ValidateResponse)(nil), // 15: trueauth.ValidateResponse
	(*SignoutResponse)(nil),  // 16: trueauth.SignoutResponse
	(*RefreshResponse)(nil),  // 17: trueauth.RefreshResponse
	(*ResetResponse)(nil),    // 18: trueauth.ResetResponse
	(*VerifyResponse)(nil),   // 19: trueauth.VerifyResponse
	(*DeleteResponse)(nil),   // 20: trueauth.DeleteResponse
	(*UserResponse)(nil),     // 21: trueauth.UserResponse
}
var file_service_trueauth_proto_depIdxs = []int32{
	0,  // 0: trueauth.TrueAuth.Welcome:input_type -> trueauth.WelcomeRequest
	1,  // 1: trueauth.TrueAuth.Health:input_type -> trueauth.HealthRequest
	2,  // 2: trueauth.TrueAuth.Signup:input_type -> trueauth.SignupRequest
	3,  // 3: trueauth.TrueAuth.Signin:input_type -> trueauth.SigninRequest
	4,  // 4: trueauth.TrueAuth.Validate:input_type -> trueauth.ValidateRequest
	5,  // 5: trueauth.TrueAuth.Signout:input_type -> trueauth.SignoutRequest
	6,  // 6: trueauth.TrueAuth.Refresh:input_type -> trueauth.RefreshRequest
	7,  // 7: trueauth.TrueAuth.Reset:input_type -> trueauth.ResetRequest
	8,  // 8: trueauth.TrueAuth.Verify:input_type -> trueauth.VerifyRequest
	9,  // 9: trueauth.TrueAuth.Delete:input_type -> trueauth.DeleteRequest
	10, // 10: trueauth.TrueAuth.User:input_type -> trueauth.UserRequest
	11, // 11: trueauth.TrueAuth.Welcome:output_type -> trueauth.WelcomeResponse
	12, // 12: trueauth.TrueAuth.Health:output_type -> trueauth.HealthResponse
	13, // 13: trueauth.TrueAuth.Signup:output_type -> trueauth.SignupResponse
	14, // 14: trueauth.TrueAuth.Signin:output_type -> trueauth.SigninResponse
	15, // 15: trueauth.TrueAuth.Validate:output_type -> trueauth.ValidateResponse
	16, // 16: trueauth.TrueAuth.Signout:output_type -> trueauth.SignoutResponse
	17, // 17: trueauth.TrueAuth.Refresh:output_type -> trueauth.RefreshResponse
	18, // 18: trueauth.TrueAuth.Reset:output_type -> trueauth.ResetResponse
	19, // 19: trueauth.TrueAuth.Verify:output_type -> trueauth.VerifyResponse
	20, // 20: trueauth.TrueAuth.Delete:output_type -> trueauth.DeleteResponse
	21, // 21: trueauth.TrueAuth.User:output_type -> trueauth.UserResponse
	11, // [11:22] is the sub-list for method output_type
	0,  // [0:11] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_service_trueauth_proto_init() }
func file_service_trueauth_proto_init() {
	if File_service_trueauth_proto != nil {
		return
	}
	file_rpc_delete_trueauth_proto_init()
	file_rpc_health_trueauth_proto_init()
	file_rpc_refresh_trueauth_proto_init()
	file_rpc_reset_trueauth_proto_init()
	file_rpc_signin_trueauth_proto_init()
	file_rpc_signout_trueauth_proto_init()
	file_rpc_signup_trueauth_proto_init()
	file_rpc_user_trueauth_proto_init()
	file_rpc_validate_trueauth_proto_init()
	file_rpc_verify_trueauth_proto_init()
	file_rpc_welcome_trueauth_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_trueauth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_trueauth_proto_goTypes,
		DependencyIndexes: file_service_trueauth_proto_depIdxs,
	}.Build()
	File_service_trueauth_proto = out.File
	file_service_trueauth_proto_rawDesc = nil
	file_service_trueauth_proto_goTypes = nil
	file_service_trueauth_proto_depIdxs = nil
}