// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.0
// source: biz/v1/error.proto

package v1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type CarPlayerErrorReason int32

const (
	// 为某个枚举单独设置错误码
	CarPlayerErrorReason_CAR_NOT_FOUND          CarPlayerErrorReason = 0
	CarPlayerErrorReason_BAND_REPEAT            CarPlayerErrorReason = 1
	CarPlayerErrorReason_BAND_OTHER_PLAYER_BAND CarPlayerErrorReason = 2
	CarPlayerErrorReason_UNBAND_REPEAT          CarPlayerErrorReason = 3
	CarPlayerErrorReason_UNBAND_NOT_BANDED      CarPlayerErrorReason = 4
)

// Enum value maps for CarPlayerErrorReason.
var (
	CarPlayerErrorReason_name = map[int32]string{
		0: "CAR_NOT_FOUND",
		1: "BAND_REPEAT",
		2: "BAND_OTHER_PLAYER_BAND",
		3: "UNBAND_REPEAT",
		4: "UNBAND_NOT_BANDED",
	}
	CarPlayerErrorReason_value = map[string]int32{
		"CAR_NOT_FOUND":          0,
		"BAND_REPEAT":            1,
		"BAND_OTHER_PLAYER_BAND": 2,
		"UNBAND_REPEAT":          3,
		"UNBAND_NOT_BANDED":      4,
	}
)

func (x CarPlayerErrorReason) Enum() *CarPlayerErrorReason {
	p := new(CarPlayerErrorReason)
	*p = x
	return p
}

func (x CarPlayerErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CarPlayerErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_biz_v1_error_proto_enumTypes[0].Descriptor()
}

func (CarPlayerErrorReason) Type() protoreflect.EnumType {
	return &file_biz_v1_error_proto_enumTypes[0]
}

func (x CarPlayerErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CarPlayerErrorReason.Descriptor instead.
func (CarPlayerErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_biz_v1_error_proto_rawDescGZIP(), []int{0}
}

var File_biz_v1_error_proto protoreflect.FileDescriptor

var file_biz_v1_error_proto_rawDesc = []byte{
	0x0a, 0x12, 0x62, 0x69, 0x7a, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x7a, 0x2e, 0x76, 0x31,
	0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xa4, 0x01, 0x0a, 0x14, 0x43, 0x61, 0x72, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x17,
	0x0a, 0x0d, 0x43, 0x41, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10,
	0x00, 0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x15, 0x0a, 0x0b, 0x42, 0x41, 0x4e, 0x44, 0x5f,
	0x52, 0x45, 0x50, 0x45, 0x41, 0x54, 0x10, 0x01, 0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x20,
	0x0a, 0x16, 0x42, 0x41, 0x4e, 0x44, 0x5f, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x5f, 0x50, 0x4c, 0x41,
	0x59, 0x45, 0x52, 0x5f, 0x42, 0x41, 0x4e, 0x44, 0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03,
	0x12, 0x17, 0x0a, 0x0d, 0x55, 0x4e, 0x42, 0x41, 0x4e, 0x44, 0x5f, 0x52, 0x45, 0x50, 0x45, 0x41,
	0x54, 0x10, 0x03, 0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x1b, 0x0a, 0x11, 0x55, 0x4e, 0x42,
	0x41, 0x4e, 0x44, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x42, 0x41, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x04,
	0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x0f, 0x5a, 0x0d,
	0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x7a, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_biz_v1_error_proto_rawDescOnce sync.Once
	file_biz_v1_error_proto_rawDescData = file_biz_v1_error_proto_rawDesc
)

func file_biz_v1_error_proto_rawDescGZIP() []byte {
	file_biz_v1_error_proto_rawDescOnce.Do(func() {
		file_biz_v1_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_biz_v1_error_proto_rawDescData)
	})
	return file_biz_v1_error_proto_rawDescData
}

var file_biz_v1_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_biz_v1_error_proto_goTypes = []any{
	(CarPlayerErrorReason)(0), // 0: api.biz.v1.CarPlayerErrorReason
}
var file_biz_v1_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_biz_v1_error_proto_init() }
func file_biz_v1_error_proto_init() {
	if File_biz_v1_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_biz_v1_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_biz_v1_error_proto_goTypes,
		DependencyIndexes: file_biz_v1_error_proto_depIdxs,
		EnumInfos:         file_biz_v1_error_proto_enumTypes,
	}.Build()
	File_biz_v1_error_proto = out.File
	file_biz_v1_error_proto_rawDesc = nil
	file_biz_v1_error_proto_goTypes = nil
	file_biz_v1_error_proto_depIdxs = nil
}
