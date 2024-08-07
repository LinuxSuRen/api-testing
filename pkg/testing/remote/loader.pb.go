// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
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

type Configs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Config `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *Configs) Reset() {
	*x = Configs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_testing_remote_loader_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configs) ProtoMessage() {}

func (x *Configs) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_testing_remote_loader_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configs.ProtoReflect.Descriptor instead.
func (*Configs) Descriptor() ([]byte, []int) {
	return file_pkg_testing_remote_loader_proto_rawDescGZIP(), []int{2}
}

func (x *Configs) GetData() []*Config {
	if x != nil {
		return x.Data
	}
	return nil
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value       string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_testing_remote_loader_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_testing_remote_loader_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_pkg_testing_remote_loader_proto_rawDescGZIP(), []int{3}
}

func (x *Config) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Config) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Config) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
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
	0x22, 0x2d, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x22, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x54, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0xc7, 0x05, 0x0a, 0x06, 0x4c, 0x6f, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x34, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74,
	0x65, 0x12, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x12, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75,
	0x69, 0x74, 0x65, 0x73, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x0d, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a,
	0x0c, 0x47, 0x65, 0x74, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x11, 0x2e,
	0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65,
	0x1a, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75,
	0x69, 0x74, 0x65, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x11, 0x2e, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x22, 0x00,
	0x12, 0x35, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75,
	0x69, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x54,
	0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x73, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x11, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x73, 0x22, 0x00,
	0x12, 0x33, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61,
	0x73, 0x65, 0x12, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74,
	0x43, 0x61, 0x73, 0x65, 0x1a, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x65, 0x73, 0x74,
	0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65,
	0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x1a, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x1a, 0x10,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x33, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74,
	0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x54, 0x65,
	0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x1a, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x06, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x12, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x17, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x05, 0x50,
	0x50, 0x72, 0x6f, 0x66, 0x12, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x50, 0x50,
	0x72, 0x6f, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x50, 0x50, 0x72, 0x6f, 0x66, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x32,
	0x96, 0x02, 0x0a, 0x0d, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x2d, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x0e,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x1a, 0x0e,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x00,
	0x12, 0x2e, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x12, 0x0d,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x22, 0x00,
	0x12, 0x36, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x12, 0x0e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x0e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00,
	0x12, 0x36, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x12, 0x0e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x32, 0x9e, 0x02, 0x0a, 0x0d, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x0a, 0x47, 0x65,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x0d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x09, 0x47, 0x65,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x12, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x0e, 0x2e, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x00, 0x12, 0x36, 0x0a,
	0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x0e, 0x2e,
	0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x14, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x0e, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x3a, 0x0a,
	0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x12, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x73, 0x75, 0x72,
	0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_pkg_testing_remote_loader_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_testing_remote_loader_proto_goTypes = []interface{}{
	(*TestSuites)(nil),             // 0: remote.TestSuites
	(*TestSuite)(nil),              // 1: remote.TestSuite
	(*Configs)(nil),                // 2: remote.Configs
	(*Config)(nil),                 // 3: remote.Config
	(*server.Pair)(nil),            // 4: server.Pair
	(*server.APISpec)(nil),         // 5: server.APISpec
	(*server.TestCase)(nil),        // 6: server.TestCase
	(*server.Empty)(nil),           // 7: server.Empty
	(*server.PProfRequest)(nil),    // 8: server.PProfRequest
	(*server.Secret)(nil),          // 9: server.Secret
	(*server.SimpleName)(nil),      // 10: server.SimpleName
	(*server.TestCases)(nil),       // 11: server.TestCases
	(*server.Version)(nil),         // 12: server.Version
	(*server.ExtensionStatus)(nil), // 13: server.ExtensionStatus
	(*server.PProfData)(nil),       // 14: server.PProfData
	(*server.Secrets)(nil),         // 15: server.Secrets
	(*server.CommonResult)(nil),    // 16: server.CommonResult
}
var file_pkg_testing_remote_loader_proto_depIdxs = []int32{
	1,  // 0: remote.TestSuites.data:type_name -> remote.TestSuite
	4,  // 1: remote.TestSuite.param:type_name -> server.Pair
	5,  // 2: remote.TestSuite.spec:type_name -> server.APISpec
	6,  // 3: remote.TestSuite.items:type_name -> server.TestCase
	3,  // 4: remote.Configs.data:type_name -> remote.Config
	7,  // 5: remote.Loader.ListTestSuite:input_type -> server.Empty
	1,  // 6: remote.Loader.CreateTestSuite:input_type -> remote.TestSuite
	1,  // 7: remote.Loader.GetTestSuite:input_type -> remote.TestSuite
	1,  // 8: remote.Loader.UpdateTestSuite:input_type -> remote.TestSuite
	1,  // 9: remote.Loader.DeleteTestSuite:input_type -> remote.TestSuite
	1,  // 10: remote.Loader.ListTestCases:input_type -> remote.TestSuite
	6,  // 11: remote.Loader.CreateTestCase:input_type -> server.TestCase
	6,  // 12: remote.Loader.GetTestCase:input_type -> server.TestCase
	6,  // 13: remote.Loader.UpdateTestCase:input_type -> server.TestCase
	6,  // 14: remote.Loader.DeleteTestCase:input_type -> server.TestCase
	7,  // 15: remote.Loader.GetVersion:input_type -> server.Empty
	7,  // 16: remote.Loader.Verify:input_type -> server.Empty
	8,  // 17: remote.Loader.PProf:input_type -> server.PProfRequest
	9,  // 18: remote.SecretService.GetSecret:input_type -> server.Secret
	7,  // 19: remote.SecretService.GetSecrets:input_type -> server.Empty
	9,  // 20: remote.SecretService.CreateSecret:input_type -> server.Secret
	9,  // 21: remote.SecretService.DeleteSecret:input_type -> server.Secret
	9,  // 22: remote.SecretService.UpdateSecret:input_type -> server.Secret
	7,  // 23: remote.ConfigService.GetConfigs:input_type -> server.Empty
	10, // 24: remote.ConfigService.GetConfig:input_type -> server.SimpleName
	3,  // 25: remote.ConfigService.CreateConfig:input_type -> remote.Config
	3,  // 26: remote.ConfigService.UpdateConfig:input_type -> remote.Config
	10, // 27: remote.ConfigService.DeleteConfig:input_type -> server.SimpleName
	0,  // 28: remote.Loader.ListTestSuite:output_type -> remote.TestSuites
	7,  // 29: remote.Loader.CreateTestSuite:output_type -> server.Empty
	1,  // 30: remote.Loader.GetTestSuite:output_type -> remote.TestSuite
	1,  // 31: remote.Loader.UpdateTestSuite:output_type -> remote.TestSuite
	7,  // 32: remote.Loader.DeleteTestSuite:output_type -> server.Empty
	11, // 33: remote.Loader.ListTestCases:output_type -> server.TestCases
	7,  // 34: remote.Loader.CreateTestCase:output_type -> server.Empty
	6,  // 35: remote.Loader.GetTestCase:output_type -> server.TestCase
	6,  // 36: remote.Loader.UpdateTestCase:output_type -> server.TestCase
	7,  // 37: remote.Loader.DeleteTestCase:output_type -> server.Empty
	12, // 38: remote.Loader.GetVersion:output_type -> server.Version
	13, // 39: remote.Loader.Verify:output_type -> server.ExtensionStatus
	14, // 40: remote.Loader.PProf:output_type -> server.PProfData
	9,  // 41: remote.SecretService.GetSecret:output_type -> server.Secret
	15, // 42: remote.SecretService.GetSecrets:output_type -> server.Secrets
	16, // 43: remote.SecretService.CreateSecret:output_type -> server.CommonResult
	16, // 44: remote.SecretService.DeleteSecret:output_type -> server.CommonResult
	16, // 45: remote.SecretService.UpdateSecret:output_type -> server.CommonResult
	2,  // 46: remote.ConfigService.GetConfigs:output_type -> remote.Configs
	3,  // 47: remote.ConfigService.GetConfig:output_type -> remote.Config
	16, // 48: remote.ConfigService.CreateConfig:output_type -> server.CommonResult
	16, // 49: remote.ConfigService.UpdateConfig:output_type -> server.CommonResult
	16, // 50: remote.ConfigService.DeleteConfig:output_type -> server.CommonResult
	28, // [28:51] is the sub-list for method output_type
	5,  // [5:28] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
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
		file_pkg_testing_remote_loader_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configs); i {
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
		file_pkg_testing_remote_loader_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   3,
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
