// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: pkg/proto/based32.proto

package proto

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

type Error int32

const (
	Error_ZERO_LENGTH                   Error = 0
	Error_CHECK_FAILED                  Error = 1
	Error_NIL_SLICE                     Error = 2
	Error_CHECK_TOO_SHORT               Error = 3
	Error_INCORRECT_HUMAN_READABLE_PART Error = 4
)

// Enum value maps for Error.
var (
	Error_name = map[int32]string{
		0: "ZERO_LENGTH",
		1: "CHECK_FAILED",
		2: "NIL_SLICE",
		3: "CHECK_TOO_SHORT",
		4: "INCORRECT_HUMAN_READABLE_PART",
	}
	Error_value = map[string]int32{
		"ZERO_LENGTH":                   0,
		"CHECK_FAILED":                  1,
		"NIL_SLICE":                     2,
		"CHECK_TOO_SHORT":               3,
		"INCORRECT_HUMAN_READABLE_PART": 4,
	}
)

func (x Error) Enum() *Error {
	p := new(Error)
	*p = x
	return p
}

func (x Error) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Error) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_proto_based32_proto_enumTypes[0].Descriptor()
}

func (Error) Type() protoreflect.EnumType {
	return &file_pkg_proto_based32_proto_enumTypes[0]
}

func (x Error) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Error.Descriptor instead.
func (Error) EnumDescriptor() ([]byte, []int) {
	return file_pkg_proto_based32_proto_rawDescGZIP(), []int{0}
}

type EncodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *EncodeRequest) Reset() {
	*x = EncodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_based32_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncodeRequest) ProtoMessage() {}

func (x *EncodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_based32_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncodeRequest.ProtoReflect.Descriptor instead.
func (*EncodeRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_based32_proto_rawDescGZIP(), []int{0}
}

func (x *EncodeRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type EncodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Encoded:
	//	*EncodeResponse_EncodedString
	//	*EncodeResponse_Error
	Encoded isEncodeResponse_Encoded `protobuf_oneof:"Encoded"`
}

func (x *EncodeResponse) Reset() {
	*x = EncodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_based32_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncodeResponse) ProtoMessage() {}

func (x *EncodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_based32_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncodeResponse.ProtoReflect.Descriptor instead.
func (*EncodeResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_based32_proto_rawDescGZIP(), []int{1}
}

func (m *EncodeResponse) GetEncoded() isEncodeResponse_Encoded {
	if m != nil {
		return m.Encoded
	}
	return nil
}

func (x *EncodeResponse) GetEncodedString() string {
	if x, ok := x.GetEncoded().(*EncodeResponse_EncodedString); ok {
		return x.EncodedString
	}
	return ""
}

func (x *EncodeResponse) GetError() Error {
	if x, ok := x.GetEncoded().(*EncodeResponse_Error); ok {
		return x.Error
	}
	return Error_ZERO_LENGTH
}

type isEncodeResponse_Encoded interface {
	isEncodeResponse_Encoded()
}

type EncodeResponse_EncodedString struct {
	EncodedString string `protobuf:"bytes,1,opt,name=EncodedString,proto3,oneof"`
}

type EncodeResponse_Error struct {
	Error Error `protobuf:"varint,2,opt,name=Error,proto3,enum=codec.Error,oneof"`
}

func (*EncodeResponse_EncodedString) isEncodeResponse_Encoded() {}

func (*EncodeResponse_Error) isEncodeResponse_Encoded() {}

type DecodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncodedString string `protobuf:"bytes,1,opt,name=EncodedString,proto3" json:"EncodedString,omitempty"`
}

func (x *DecodeRequest) Reset() {
	*x = DecodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_based32_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecodeRequest) ProtoMessage() {}

func (x *DecodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_based32_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecodeRequest.ProtoReflect.Descriptor instead.
func (*DecodeRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_based32_proto_rawDescGZIP(), []int{2}
}

func (x *DecodeRequest) GetEncodedString() string {
	if x != nil {
		return x.EncodedString
	}
	return ""
}

type DecodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Decoded:
	//	*DecodeResponse_Data
	//	*DecodeResponse_Error
	Decoded isDecodeResponse_Decoded `protobuf_oneof:"Decoded"`
}

func (x *DecodeResponse) Reset() {
	*x = DecodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_based32_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecodeResponse) ProtoMessage() {}

func (x *DecodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_based32_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecodeResponse.ProtoReflect.Descriptor instead.
func (*DecodeResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_based32_proto_rawDescGZIP(), []int{3}
}

func (m *DecodeResponse) GetDecoded() isDecodeResponse_Decoded {
	if m != nil {
		return m.Decoded
	}
	return nil
}

func (x *DecodeResponse) GetData() []byte {
	if x, ok := x.GetDecoded().(*DecodeResponse_Data); ok {
		return x.Data
	}
	return nil
}

func (x *DecodeResponse) GetError() Error {
	if x, ok := x.GetDecoded().(*DecodeResponse_Error); ok {
		return x.Error
	}
	return Error_ZERO_LENGTH
}

type isDecodeResponse_Decoded interface {
	isDecodeResponse_Decoded()
}

type DecodeResponse_Data struct {
	Data []byte `protobuf:"bytes,1,opt,name=Data,proto3,oneof"`
}

type DecodeResponse_Error struct {
	Error Error `protobuf:"varint,2,opt,name=Error,proto3,enum=codec.Error,oneof"`
}

func (*DecodeResponse_Data) isDecodeResponse_Decoded() {}

func (*DecodeResponse_Error) isDecodeResponse_Decoded() {}

var File_pkg_proto_based32_proto protoreflect.FileDescriptor

var file_pkg_proto_based32_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x61, 0x73, 0x65,
	0x64, 0x33, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x63,
	0x22, 0x23, 0x0a, 0x0d, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0x69, 0x0a, 0x0e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x0d, 0x45, 0x6e, 0x63, 0x6f, 0x64,
	0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x0d, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12,
	0x24, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c,
	0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x09, 0x0a, 0x07, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64,
	0x22, 0x35, 0x0a, 0x0d, 0x44, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x24, 0x0a, 0x0d, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65,
	0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x57, 0x0a, 0x0e, 0x44, 0x65, 0x63, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x04, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x24, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c,
	0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x09, 0x0a, 0x07, 0x44, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x64,
	0x2a, 0x71, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x0f, 0x0a, 0x0b, 0x5a, 0x45, 0x52,
	0x4f, 0x5f, 0x4c, 0x45, 0x4e, 0x47, 0x54, 0x48, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x48,
	0x45, 0x43, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09,
	0x4e, 0x49, 0x4c, 0x5f, 0x53, 0x4c, 0x49, 0x43, 0x45, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x43,
	0x48, 0x45, 0x43, 0x4b, 0x5f, 0x54, 0x4f, 0x4f, 0x5f, 0x53, 0x48, 0x4f, 0x52, 0x54, 0x10, 0x03,
	0x12, 0x21, 0x0a, 0x1d, 0x49, 0x4e, 0x43, 0x4f, 0x52, 0x52, 0x45, 0x43, 0x54, 0x5f, 0x48, 0x55,
	0x4d, 0x41, 0x4e, 0x5f, 0x52, 0x45, 0x41, 0x44, 0x41, 0x42, 0x4c, 0x45, 0x5f, 0x50, 0x41, 0x52,
	0x54, 0x10, 0x04, 0x32, 0x83, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x06, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x2e,
	0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x45, 0x6e, 0x63, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x39,
	0x0a, 0x06, 0x44, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63,
	0x2e, 0x44, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x2e, 0x44, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x65, 0x72, 0x61,
	0x6c, 0x6c, 0x2f, 0x6b, 0x69, 0x74, 0x63, 0x68, 0x65, 0x6e, 0x73, 0x69, 0x6e, 0x6b, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_based32_proto_rawDescOnce sync.Once
	file_pkg_proto_based32_proto_rawDescData = file_pkg_proto_based32_proto_rawDesc
)

func file_pkg_proto_based32_proto_rawDescGZIP() []byte {
	file_pkg_proto_based32_proto_rawDescOnce.Do(func() {
		file_pkg_proto_based32_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_based32_proto_rawDescData)
	})
	return file_pkg_proto_based32_proto_rawDescData
}

var file_pkg_proto_based32_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_proto_based32_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_proto_based32_proto_goTypes = []interface{}{
	(Error)(0),             // 0: codec.Error
	(*EncodeRequest)(nil),  // 1: codec.EncodeRequest
	(*EncodeResponse)(nil), // 2: codec.EncodeResponse
	(*DecodeRequest)(nil),  // 3: codec.DecodeRequest
	(*DecodeResponse)(nil), // 4: codec.DecodeResponse
}
var file_pkg_proto_based32_proto_depIdxs = []int32{
	0, // 0: codec.EncodeResponse.Error:type_name -> codec.Error
	0, // 1: codec.DecodeResponse.Error:type_name -> codec.Error
	1, // 2: codec.Transcriber.Encode:input_type -> codec.EncodeRequest
	3, // 3: codec.Transcriber.Decode:input_type -> codec.DecodeRequest
	2, // 4: codec.Transcriber.Encode:output_type -> codec.EncodeResponse
	4, // 5: codec.Transcriber.Decode:output_type -> codec.DecodeResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_proto_based32_proto_init() }
func file_pkg_proto_based32_proto_init() {
	if File_pkg_proto_based32_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_based32_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncodeRequest); i {
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
		file_pkg_proto_based32_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncodeResponse); i {
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
		file_pkg_proto_based32_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DecodeRequest); i {
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
		file_pkg_proto_based32_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DecodeResponse); i {
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
	file_pkg_proto_based32_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*EncodeResponse_EncodedString)(nil),
		(*EncodeResponse_Error)(nil),
	}
	file_pkg_proto_based32_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*DecodeResponse_Data)(nil),
		(*DecodeResponse_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_proto_based32_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_based32_proto_goTypes,
		DependencyIndexes: file_pkg_proto_based32_proto_depIdxs,
		EnumInfos:         file_pkg_proto_based32_proto_enumTypes,
		MessageInfos:      file_pkg_proto_based32_proto_msgTypes,
	}.Build()
	File_pkg_proto_based32_proto = out.File
	file_pkg_proto_based32_proto_rawDesc = nil
	file_pkg_proto_based32_proto_goTypes = nil
	file_pkg_proto_based32_proto_depIdxs = nil
}
