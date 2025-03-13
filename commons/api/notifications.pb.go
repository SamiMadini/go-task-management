// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: api/notifications.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NotificationType int32

const (
	NotificationType_IN_APP NotificationType = 0
	NotificationType_EMAIL  NotificationType = 1
)

// Enum value maps for NotificationType.
var (
	NotificationType_name = map[int32]string{
		0: "IN_APP",
		1: "EMAIL",
	}
	NotificationType_value = map[string]int32{
		"IN_APP": 0,
		"EMAIL":  1,
	}
)

func (x NotificationType) Enum() *NotificationType {
	p := new(NotificationType)
	*p = x
	return p
}

func (x NotificationType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NotificationType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_notifications_proto_enumTypes[0].Descriptor()
}

func (NotificationType) Type() protoreflect.EnumType {
	return &file_api_notifications_proto_enumTypes[0]
}

func (x NotificationType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NotificationType.Descriptor instead.
func (NotificationType) EnumDescriptor() ([]byte, []int) {
	return file_api_notifications_proto_rawDescGZIP(), []int{0}
}

type SendNotificationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskId        string                 `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	Types         []NotificationType     `protobuf:"varint,2,rep,packed,name=types,proto3,enum=api.NotificationType" json:"types,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendNotificationRequest) Reset() {
	*x = SendNotificationRequest{}
	mi := &file_api_notifications_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendNotificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendNotificationRequest) ProtoMessage() {}

func (x *SendNotificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_notifications_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendNotificationRequest.ProtoReflect.Descriptor instead.
func (*SendNotificationRequest) Descriptor() ([]byte, []int) {
	return file_api_notifications_proto_rawDescGZIP(), []int{0}
}

func (x *SendNotificationRequest) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *SendNotificationRequest) GetTypes() []NotificationType {
	if x != nil {
		return x.Types
	}
	return nil
}

type SendNotificationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ack           string                 `protobuf:"bytes,1,opt,name=ack,proto3" json:"ack,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendNotificationResponse) Reset() {
	*x = SendNotificationResponse{}
	mi := &file_api_notifications_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendNotificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendNotificationResponse) ProtoMessage() {}

func (x *SendNotificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_notifications_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendNotificationResponse.ProtoReflect.Descriptor instead.
func (*SendNotificationResponse) Descriptor() ([]byte, []int) {
	return file_api_notifications_proto_rawDescGZIP(), []int{1}
}

func (x *SendNotificationResponse) GetAck() string {
	if x != nil {
		return x.Ack
	}
	return ""
}

var File_api_notifications_proto protoreflect.FileDescriptor

var file_api_notifications_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x22, 0x5e,
	0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x73,
	0x6b, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49,
	0x64, 0x12, 0x2b, 0x0a, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e,
	0x32, 0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x22, 0x2c,
	0x0a, 0x18, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x63, 0x6b, 0x2a, 0x29, 0x0a, 0x10,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x0a, 0x0a, 0x06, 0x49, 0x4e, 0x5f, 0x41, 0x50, 0x50, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05,
	0x45, 0x4d, 0x41, 0x49, 0x4c, 0x10, 0x01, 0x32, 0x68, 0x0a, 0x13, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x51,
	0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x25, 0x5a, 0x23, 0x73, 0x61, 0x6d, 0x61, 0x2f, 0x67, 0x6f, 0x2d, 0x74, 0x61, 0x73,
	0x6b, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_notifications_proto_rawDescOnce sync.Once
	file_api_notifications_proto_rawDescData []byte
)

func file_api_notifications_proto_rawDescGZIP() []byte {
	file_api_notifications_proto_rawDescOnce.Do(func() {
		file_api_notifications_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_notifications_proto_rawDesc), len(file_api_notifications_proto_rawDesc)))
	})
	return file_api_notifications_proto_rawDescData
}

var file_api_notifications_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_notifications_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_notifications_proto_goTypes = []any{
	(NotificationType)(0),            // 0: api.NotificationType
	(*SendNotificationRequest)(nil),  // 1: api.SendNotificationRequest
	(*SendNotificationResponse)(nil), // 2: api.SendNotificationResponse
}
var file_api_notifications_proto_depIdxs = []int32{
	0, // 0: api.SendNotificationRequest.types:type_name -> api.NotificationType
	1, // 1: api.NotificationService.SendNotification:input_type -> api.SendNotificationRequest
	2, // 2: api.NotificationService.SendNotification:output_type -> api.SendNotificationResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_notifications_proto_init() }
func file_api_notifications_proto_init() {
	if File_api_notifications_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_notifications_proto_rawDesc), len(file_api_notifications_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_notifications_proto_goTypes,
		DependencyIndexes: file_api_notifications_proto_depIdxs,
		EnumInfos:         file_api_notifications_proto_enumTypes,
		MessageInfos:      file_api_notifications_proto_msgTypes,
	}.Build()
	File_api_notifications_proto = out.File
	file_api_notifications_proto_goTypes = nil
	file_api_notifications_proto_depIdxs = nil
}
