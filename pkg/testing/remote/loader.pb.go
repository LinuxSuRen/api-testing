// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.2
// source: pkg/testing/remote/loader.proto

package remote

import (
	server "github.com/linuxsuren/api-testing/pkg/server"
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

type TestSuites struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*TestSuite `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *TestSuites) Reset() {
	*x = TestSuites{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_testing_remote_loader_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestSuites) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestSuites) ProtoMessage() {}

func (x *TestSuites) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_testing_remote_loader_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestSuites.ProtoReflect.Descriptor instead.
func (*TestSuites) Descriptor() ([]byte, []int) {
	return file_pkg_testing_remote_loader_proto_rawDescGZIP(), []int{0}
}

func (x *TestSuites) GetData() []*TestSuite {
	if x != nil {
		return x.Data
	}
	return nil
}

type TestSuite struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Api   string             `protobuf:"bytes,2,opt,name=api,proto3" json:"api,omitempty"`
	Param []*server.Pair     `protobuf:"bytes,3,rep,name=param,proto3" json:"param,omitempty"`
	Spec  *server.APISpec    `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Items []*server.TestCase `protobuf:"bytes,5,rep,name=items,proto3" json:"items,omitempty"`
	Full  bool               `protobuf:"varint,6,opt,name=full,proto3" json:"full,omitempty"`
}

func (x *TestSuite) Reset() {
	*x = TestSuite{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_testing_remote_loader_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestSuite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestSuite) ProtoMessage() {}

func (x *TestSuite) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_testing_remote_loader_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestSuite.ProtoReflect.Descriptor instead.
func (*TestSuite) Descriptor() ([]byte, []int) {
	return file_pkg_testing_remote_loader_proto_rawDescGZIP(), []int{1}
}

func (x *TestSuite) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TestSuite) GetApi() string {
	if x != nil {
		return x.Api
	}
	return ""
}

func (x *TestSuite) GetParam() []*server.Pair {
	if x != nil {
		return x.Param
	}
	return nil
}

func (x *TestSuite) GetSpec() *server.APISpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *TestSuite) GetItems() []*server.TestCase {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *TestSuite) GetFull() bool {
	if x != nil {
		return x.Full
	}
	return false
}

var File_pkg_testing_remote_loader_proto protoreflect.FileDescriptor

var file_pkg_testing_remote_loader_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x2f, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x1a, 0x17, 0x70, 0x6b, 0x67, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x33, 0x0a, 0x0a, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x73,
	0x12, 0x25, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74,
	0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xb6, 0x01, 0x0a, 0x09, 0x54, 0x65, 0x73, 0x74,
	0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x70, 0x69,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x69, 0x12, 0x22, 0x0a, 0x05, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x50, 0x61, 0x69, 0x72, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12,
	0x23, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x41, 0x50, 0x49, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04,
	0x73, 0x70, 0x65, 0x63, 0x12, 0x26, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x43, 0x61, 0x73, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x66, 0x75, 0x6c, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x66, 0x75, 0x6c, 0x6c,
	0x32, 0xe0, 0x04, 0x0a, 0x06, 0x4c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x12, 0x34, 0x0a, 0x0d, 0x4c,
	0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x0d, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x73, 0x22,
	0x00, 0x12, 0x35, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x53,
	0x75, 0x69, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65,
	0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x54,
	0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x11, 0x2e, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x22, 0x00,
	0x12, 0x39, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75,
	0x69, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x0f, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x11,
	0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74,
	0x65, 0x1a, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x37, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61,
	0x73, 0x65, 0x73, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x73, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x0e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x1a,
	0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x33, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12,
	0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73,
	0x65, 0x1a, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43,
	0x61, 0x73, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x1a, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x22, 0x00, 0x12, 0x33, 0x0a,
	0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12,
	0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73,
	0x65, 0x1a, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x2f, 0x0a, 0x06, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12, 0x0d, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x00, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x73, 0x75, 0x72, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69,
	0x2d, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x67, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_testing_remote_loader_proto_rawDescOnce sync.Once
	file_pkg_testing_remote_loader_proto_rawDescData = file_pkg_testing_remote_loader_proto_rawDesc
)

func file_pkg_testing_remote_loader_proto_rawDescGZIP() []byte {
	file_pkg_testing_remote_loader_proto_rawDescOnce.Do(func() {
		file_pkg_testing_remote_loader_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_testing_remote_loader_proto_rawDescData)
	})
	return file_pkg_testing_remote_loader_proto_rawDescData
}

var file_pkg_testing_remote_loader_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_testing_remote_loader_proto_goTypes = []interface{}{
	(*TestSuites)(nil),          // 0: remote.TestSuites
	(*TestSuite)(nil),           // 1: remote.TestSuite
	(*server.Pair)(nil),         // 2: server.Pair
	(*server.APISpec)(nil),      // 3: server.APISpec
	(*server.TestCase)(nil),     // 4: server.TestCase
	(*server.Empty)(nil),        // 5: server.Empty
	(*server.TestCases)(nil),    // 6: server.TestCases
	(*server.CommonResult)(nil), // 7: server.CommonResult
}
var file_pkg_testing_remote_loader_proto_depIdxs = []int32{
	1,  // 0: remote.TestSuites.data:type_name -> remote.TestSuite
	2,  // 1: remote.TestSuite.param:type_name -> server.Pair
	3,  // 2: remote.TestSuite.spec:type_name -> server.APISpec
	4,  // 3: remote.TestSuite.items:type_name -> server.TestCase
	5,  // 4: remote.Loader.ListTestSuite:input_type -> server.Empty
	1,  // 5: remote.Loader.CreateTestSuite:input_type -> remote.TestSuite
	1,  // 6: remote.Loader.GetTestSuite:input_type -> remote.TestSuite
	1,  // 7: remote.Loader.UpdateTestSuite:input_type -> remote.TestSuite
	1,  // 8: remote.Loader.DeleteTestSuite:input_type -> remote.TestSuite
	1,  // 9: remote.Loader.ListTestCases:input_type -> remote.TestSuite
	4,  // 10: remote.Loader.CreateTestCase:input_type -> server.TestCase
	4,  // 11: remote.Loader.GetTestCase:input_type -> server.TestCase
	4,  // 12: remote.Loader.UpdateTestCase:input_type -> server.TestCase
	4,  // 13: remote.Loader.DeleteTestCase:input_type -> server.TestCase
	5,  // 14: remote.Loader.Verify:input_type -> server.Empty
	0,  // 15: remote.Loader.ListTestSuite:output_type -> remote.TestSuites
	5,  // 16: remote.Loader.CreateTestSuite:output_type -> server.Empty
	1,  // 17: remote.Loader.GetTestSuite:output_type -> remote.TestSuite
	1,  // 18: remote.Loader.UpdateTestSuite:output_type -> remote.TestSuite
	5,  // 19: remote.Loader.DeleteTestSuite:output_type -> server.Empty
	6,  // 20: remote.Loader.ListTestCases:output_type -> server.TestCases
	5,  // 21: remote.Loader.CreateTestCase:output_type -> server.Empty
	4,  // 22: remote.Loader.GetTestCase:output_type -> server.TestCase
	4,  // 23: remote.Loader.UpdateTestCase:output_type -> server.TestCase
	5,  // 24: remote.Loader.DeleteTestCase:output_type -> server.Empty
	7,  // 25: remote.Loader.Verify:output_type -> server.CommonResult
	15, // [15:26] is the sub-list for method output_type
	4,  // [4:15] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_pkg_testing_remote_loader_proto_init() }
func file_pkg_testing_remote_loader_proto_init() {
	if File_pkg_testing_remote_loader_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_testing_remote_loader_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestSuites); i {
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
		file_pkg_testing_remote_loader_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestSuite); i {
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
			RawDescriptor: file_pkg_testing_remote_loader_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_testing_remote_loader_proto_goTypes,
		DependencyIndexes: file_pkg_testing_remote_loader_proto_depIdxs,
		MessageInfos:      file_pkg_testing_remote_loader_proto_msgTypes,
	}.Build()
	File_pkg_testing_remote_loader_proto = out.File
	file_pkg_testing_remote_loader_proto_rawDesc = nil
	file_pkg_testing_remote_loader_proto_goTypes = nil
	file_pkg_testing_remote_loader_proto_depIdxs = nil
}
