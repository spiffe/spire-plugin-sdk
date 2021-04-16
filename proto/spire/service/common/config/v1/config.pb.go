// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0
// source: spire/service/common/config/v1/config.proto

package configv1

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

type ConfigureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. Core SPIRE configuration.
	CoreConfiguration *CoreConfiguration `protobuf:"bytes,1,opt,name=core_configuration,json=coreConfiguration,proto3" json:"core_configuration,omitempty"`
	// Required. HCL encoded plugin configuration.
	HclConfiguration string `protobuf:"bytes,2,opt,name=hcl_configuration,json=hclConfiguration,proto3" json:"hcl_configuration,omitempty"`
}

func (x *ConfigureRequest) Reset() {
	*x = ConfigureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spire_service_common_config_v1_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureRequest) ProtoMessage() {}

func (x *ConfigureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_spire_service_common_config_v1_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigureRequest.ProtoReflect.Descriptor instead.
func (*ConfigureRequest) Descriptor() ([]byte, []int) {
	return file_spire_service_common_config_v1_config_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigureRequest) GetCoreConfiguration() *CoreConfiguration {
	if x != nil {
		return x.CoreConfiguration
	}
	return nil
}

func (x *ConfigureRequest) GetHclConfiguration() string {
	if x != nil {
		return x.HclConfiguration
	}
	return ""
}

type ConfigureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConfigureResponse) Reset() {
	*x = ConfigureResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spire_service_common_config_v1_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureResponse) ProtoMessage() {}

func (x *ConfigureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_spire_service_common_config_v1_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigureResponse.ProtoReflect.Descriptor instead.
func (*ConfigureResponse) Descriptor() ([]byte, []int) {
	return file_spire_service_common_config_v1_config_proto_rawDescGZIP(), []int{1}
}

type CoreConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. The trust domain name SPIRE is configured with (e.g.
	// "example.org").
	TrustDomain string `protobuf:"bytes,1,opt,name=trust_domain,json=trustDomain,proto3" json:"trust_domain,omitempty"`
}

func (x *CoreConfiguration) Reset() {
	*x = CoreConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spire_service_common_config_v1_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoreConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoreConfiguration) ProtoMessage() {}

func (x *CoreConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_spire_service_common_config_v1_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoreConfiguration.ProtoReflect.Descriptor instead.
func (*CoreConfiguration) Descriptor() ([]byte, []int) {
	return file_spire_service_common_config_v1_config_proto_rawDescGZIP(), []int{2}
}

func (x *CoreConfiguration) GetTrustDomain() string {
	if x != nil {
		return x.TrustDomain
	}
	return ""
}

var File_spire_service_common_config_v1_config_proto protoreflect.FileDescriptor

var file_spire_service_common_config_v1_config_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x73,
	0x70, 0x69, 0x72, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x22, 0xa1, 0x01,
	0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x60, 0x0a, 0x12, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31,
	0x2e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x6f, 0x72, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x11, 0x63, 0x6f, 0x72, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2b, 0x0a, 0x11, 0x68, 0x63, 0x6c, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x10, 0x68, 0x63, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x13, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x36, 0x0a, 0x11, 0x43, 0x6f, 0x72, 0x65, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x74,
	0x72, 0x75, 0x73, 0x74, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x74, 0x72, 0x75, 0x73, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x32, 0x7a,
	0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x70, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x65, 0x12, 0x30, 0x2e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x52, 0x5a, 0x50, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x70, 0x69, 0x66, 0x66, 0x65, 0x2f,
	0x73, 0x70, 0x69, 0x72, 0x65, 0x2d, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2d, 0x73, 0x64, 0x6b,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x70, 0x69, 0x72, 0x65, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_spire_service_common_config_v1_config_proto_rawDescOnce sync.Once
	file_spire_service_common_config_v1_config_proto_rawDescData = file_spire_service_common_config_v1_config_proto_rawDesc
)

func file_spire_service_common_config_v1_config_proto_rawDescGZIP() []byte {
	file_spire_service_common_config_v1_config_proto_rawDescOnce.Do(func() {
		file_spire_service_common_config_v1_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_spire_service_common_config_v1_config_proto_rawDescData)
	})
	return file_spire_service_common_config_v1_config_proto_rawDescData
}

var file_spire_service_common_config_v1_config_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_spire_service_common_config_v1_config_proto_goTypes = []interface{}{
	(*ConfigureRequest)(nil),  // 0: spire.service.common.config.v1.ConfigureRequest
	(*ConfigureResponse)(nil), // 1: spire.service.common.config.v1.ConfigureResponse
	(*CoreConfiguration)(nil), // 2: spire.service.common.config.v1.CoreConfiguration
}
var file_spire_service_common_config_v1_config_proto_depIdxs = []int32{
	2, // 0: spire.service.common.config.v1.ConfigureRequest.core_configuration:type_name -> spire.service.common.config.v1.CoreConfiguration
	0, // 1: spire.service.common.config.v1.Config.Configure:input_type -> spire.service.common.config.v1.ConfigureRequest
	1, // 2: spire.service.common.config.v1.Config.Configure:output_type -> spire.service.common.config.v1.ConfigureResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_spire_service_common_config_v1_config_proto_init() }
func file_spire_service_common_config_v1_config_proto_init() {
	if File_spire_service_common_config_v1_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_spire_service_common_config_v1_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigureRequest); i {
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
		file_spire_service_common_config_v1_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigureResponse); i {
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
		file_spire_service_common_config_v1_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoreConfiguration); i {
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
			RawDescriptor: file_spire_service_common_config_v1_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_spire_service_common_config_v1_config_proto_goTypes,
		DependencyIndexes: file_spire_service_common_config_v1_config_proto_depIdxs,
		MessageInfos:      file_spire_service_common_config_v1_config_proto_msgTypes,
	}.Build()
	File_spire_service_common_config_v1_config_proto = out.File
	file_spire_service_common_config_v1_config_proto_rawDesc = nil
	file_spire_service_common_config_v1_config_proto_goTypes = nil
	file_spire_service_common_config_v1_config_proto_depIdxs = nil
}
