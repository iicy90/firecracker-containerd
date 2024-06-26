// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.12.4
// source: ioproxy.proto

package ioproxy

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type StateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID     string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	ExecID string `protobuf:"bytes,2,opt,name=ExecID,proto3" json:"ExecID,omitempty"`
}

func (x *StateRequest) Reset() {
	*x = StateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ioproxy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateRequest) ProtoMessage() {}

func (x *StateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ioproxy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateRequest.ProtoReflect.Descriptor instead.
func (*StateRequest) Descriptor() ([]byte, []int) {
	return file_ioproxy_proto_rawDescGZIP(), []int{0}
}

func (x *StateRequest) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *StateRequest) GetExecID() string {
	if x != nil {
		return x.ExecID
	}
	return ""
}

type StateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsOpen bool `protobuf:"varint,1,opt,name=IsOpen,proto3" json:"IsOpen,omitempty"`
}

func (x *StateResponse) Reset() {
	*x = StateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ioproxy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateResponse) ProtoMessage() {}

func (x *StateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ioproxy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateResponse.ProtoReflect.Descriptor instead.
func (*StateResponse) Descriptor() ([]byte, []int) {
	return file_ioproxy_proto_rawDescGZIP(), []int{1}
}

func (x *StateResponse) GetIsOpen() bool {
	if x != nil {
		return x.IsOpen
	}
	return false
}

type AttachRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID         string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	ExecID     string `protobuf:"bytes,2,opt,name=ExecID,proto3" json:"ExecID,omitempty"`
	StdinPort  uint32 `protobuf:"varint,3,opt,name=StdinPort,proto3" json:"StdinPort,omitempty"`
	StdoutPort uint32 `protobuf:"varint,4,opt,name=StdoutPort,proto3" json:"StdoutPort,omitempty"`
	StderrPort uint32 `protobuf:"varint,5,opt,name=StderrPort,proto3" json:"StderrPort,omitempty"`
}

func (x *AttachRequest) Reset() {
	*x = AttachRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ioproxy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttachRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttachRequest) ProtoMessage() {}

func (x *AttachRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ioproxy_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttachRequest.ProtoReflect.Descriptor instead.
func (*AttachRequest) Descriptor() ([]byte, []int) {
	return file_ioproxy_proto_rawDescGZIP(), []int{2}
}

func (x *AttachRequest) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *AttachRequest) GetExecID() string {
	if x != nil {
		return x.ExecID
	}
	return ""
}

func (x *AttachRequest) GetStdinPort() uint32 {
	if x != nil {
		return x.StdinPort
	}
	return 0
}

func (x *AttachRequest) GetStdoutPort() uint32 {
	if x != nil {
		return x.StdoutPort
	}
	return 0
}

func (x *AttachRequest) GetStderrPort() uint32 {
	if x != nil {
		return x.StderrPort
	}
	return 0
}

var File_ioproxy_proto protoreflect.FileDescriptor

var file_ioproxy_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x69, 0x6f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x0c,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x45, 0x78, 0x65, 0x63, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x45, 0x78,
	0x65, 0x63, 0x49, 0x44, 0x22, 0x27, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x49, 0x73, 0x4f, 0x70, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x49, 0x73, 0x4f, 0x70, 0x65, 0x6e, 0x22, 0x95, 0x01,
	0x0a, 0x0d, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12,
	0x16, 0x0a, 0x06, 0x45, 0x78, 0x65, 0x63, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x45, 0x78, 0x65, 0x63, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x74, 0x64, 0x69, 0x6e,
	0x50, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x53, 0x74, 0x64, 0x69,
	0x6e, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x50,
	0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x53, 0x74, 0x64, 0x6f, 0x75,
	0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x64, 0x65, 0x72, 0x72, 0x50,
	0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x53, 0x74, 0x64, 0x65, 0x72,
	0x72, 0x50, 0x6f, 0x72, 0x74, 0x32, 0x63, 0x0a, 0x07, 0x49, 0x4f, 0x50, 0x72, 0x6f, 0x78, 0x79,
	0x12, 0x26, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0d, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x06, 0x41, 0x74, 0x74, 0x61,
	0x63, 0x68, 0x12, 0x0e, 0x2e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b,
	0x69, 0x6f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ioproxy_proto_rawDescOnce sync.Once
	file_ioproxy_proto_rawDescData = file_ioproxy_proto_rawDesc
)

func file_ioproxy_proto_rawDescGZIP() []byte {
	file_ioproxy_proto_rawDescOnce.Do(func() {
		file_ioproxy_proto_rawDescData = protoimpl.X.CompressGZIP(file_ioproxy_proto_rawDescData)
	})
	return file_ioproxy_proto_rawDescData
}

var file_ioproxy_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_ioproxy_proto_goTypes = []interface{}{
	(*StateRequest)(nil),  // 0: StateRequest
	(*StateResponse)(nil), // 1: StateResponse
	(*AttachRequest)(nil), // 2: AttachRequest
	(*empty.Empty)(nil),   // 3: google.protobuf.Empty
}
var file_ioproxy_proto_depIdxs = []int32{
	0, // 0: IOProxy.State:input_type -> StateRequest
	2, // 1: IOProxy.Attach:input_type -> AttachRequest
	1, // 2: IOProxy.State:output_type -> StateResponse
	3, // 3: IOProxy.Attach:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ioproxy_proto_init() }
func file_ioproxy_proto_init() {
	if File_ioproxy_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ioproxy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateRequest); i {
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
		file_ioproxy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateResponse); i {
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
		file_ioproxy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttachRequest); i {
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
			RawDescriptor: file_ioproxy_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ioproxy_proto_goTypes,
		DependencyIndexes: file_ioproxy_proto_depIdxs,
		MessageInfos:      file_ioproxy_proto_msgTypes,
	}.Build()
	File_ioproxy_proto = out.File
	file_ioproxy_proto_rawDesc = nil
	file_ioproxy_proto_goTypes = nil
	file_ioproxy_proto_depIdxs = nil
}
