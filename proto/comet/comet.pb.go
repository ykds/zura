// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.2
// source: proto/comet/comet.proto

package comet

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

type Op int32

const (
	Op_Empty                   Op = 0
	Op_NewMsg                  Op = 1
	Op_NewApplication          Op = 2
	Op_ApplicationHandleResult Op = 3
	Op_SynNewMsg               Op = 4
	Op_NewMsgReply             Op = 5
	Op_Heartbeat               Op = 6
	Op_HeartbeatReply          Op = 7
)

// Enum value maps for Op.
var (
	Op_name = map[int32]string{
		0: "Empty",
		1: "NewMsg",
		2: "NewApplication",
		3: "ApplicationHandleResult",
		4: "SynNewMsg",
		5: "NewMsgReply",
		6: "Heartbeat",
		7: "HeartbeatReply",
	}
	Op_value = map[string]int32{
		"Empty":                   0,
		"NewMsg":                  1,
		"NewApplication":          2,
		"ApplicationHandleResult": 3,
		"SynNewMsg":               4,
		"NewMsgReply":             5,
		"Heartbeat":               6,
		"HeartbeatReply":          7,
	}
)

func (x Op) Enum() *Op {
	p := new(Op)
	*p = x
	return p
}

func (x Op) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Op) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_comet_comet_proto_enumTypes[0].Descriptor()
}

func (Op) Type() protoreflect.EnumType {
	return &file_proto_comet_comet_proto_enumTypes[0]
}

func (x Op) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Op.Descriptor instead.
func (Op) EnumDescriptor() ([]byte, []int) {
	return file_proto_comet_comet_proto_rawDescGZIP(), []int{0}
}

type PushNotificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ToUserId []int64 `protobuf:"varint,1,rep,packed,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`
	Body     []byte  `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *PushNotificationRequest) Reset() {
	*x = PushNotificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_comet_comet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushNotificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushNotificationRequest) ProtoMessage() {}

func (x *PushNotificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_comet_comet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushNotificationRequest.ProtoReflect.Descriptor instead.
func (*PushNotificationRequest) Descriptor() ([]byte, []int) {
	return file_proto_comet_comet_proto_rawDescGZIP(), []int{0}
}

func (x *PushNotificationRequest) GetToUserId() []int64 {
	if x != nil {
		return x.ToUserId
	}
	return nil
}

func (x *PushNotificationRequest) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type PushNotificationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushNotificationResponse) Reset() {
	*x = PushNotificationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_comet_comet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushNotificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushNotificationResponse) ProtoMessage() {}

func (x *PushNotificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_comet_comet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushNotificationResponse.ProtoReflect.Descriptor instead.
func (*PushNotificationResponse) Descriptor() ([]byte, []int) {
	return file_proto_comet_comet_proto_rawDescGZIP(), []int{1}
}

var File_proto_comet_comet_proto protoreflect.FileDescriptor

var file_proto_comet_comet_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2f, 0x63, 0x6f,
	0x6d, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x22, 0x4b, 0x0a, 0x17, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x0a, 0x74,
	0x6f, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x08, 0x74, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x1a, 0x0a,
	0x18, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x8f, 0x01, 0x0a, 0x02, 0x4f, 0x70,
	0x12, 0x09, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4e,
	0x65, 0x77, 0x4d, 0x73, 0x67, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x4e, 0x65, 0x77, 0x41, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x41,
	0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x79, 0x6e, 0x4e,
	0x65, 0x77, 0x4d, 0x73, 0x67, 0x10, 0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x4e, 0x65, 0x77, 0x4d, 0x73,
	0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x05, 0x12, 0x0d, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x10, 0x06, 0x12, 0x12, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x72, 0x74,
	0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x07, 0x32, 0x5c, 0x0a, 0x05, 0x43,
	0x6f, 0x6d, 0x65, 0x74, 0x12, 0x53, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x2e, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x2e, 0x50, 0x75, 0x73, 0x68, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x6b, 0x64, 0x73, 0x2f, 0x7a, 0x75, 0x72,
	0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x3b, 0x63, 0x6f,
	0x6d, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_comet_comet_proto_rawDescOnce sync.Once
	file_proto_comet_comet_proto_rawDescData = file_proto_comet_comet_proto_rawDesc
)

func file_proto_comet_comet_proto_rawDescGZIP() []byte {
	file_proto_comet_comet_proto_rawDescOnce.Do(func() {
		file_proto_comet_comet_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_comet_comet_proto_rawDescData)
	})
	return file_proto_comet_comet_proto_rawDescData
}

var file_proto_comet_comet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_comet_comet_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_comet_comet_proto_goTypes = []interface{}{
	(Op)(0),                          // 0: comet.Op
	(*PushNotificationRequest)(nil),  // 1: comet.PushNotificationRequest
	(*PushNotificationResponse)(nil), // 2: comet.PushNotificationResponse
}
var file_proto_comet_comet_proto_depIdxs = []int32{
	1, // 0: comet.Comet.PushNotification:input_type -> comet.PushNotificationRequest
	2, // 1: comet.Comet.PushNotification:output_type -> comet.PushNotificationResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_comet_comet_proto_init() }
func file_proto_comet_comet_proto_init() {
	if File_proto_comet_comet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_comet_comet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushNotificationRequest); i {
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
		file_proto_comet_comet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushNotificationResponse); i {
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
			RawDescriptor: file_proto_comet_comet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_comet_comet_proto_goTypes,
		DependencyIndexes: file_proto_comet_comet_proto_depIdxs,
		EnumInfos:         file_proto_comet_comet_proto_enumTypes,
		MessageInfos:      file_proto_comet_comet_proto_msgTypes,
	}.Build()
	File_proto_comet_comet_proto = out.File
	file_proto_comet_comet_proto_rawDesc = nil
	file_proto_comet_comet_proto_goTypes = nil
	file_proto_comet_comet_proto_depIdxs = nil
}
