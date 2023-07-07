// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.2
// source: remote/loader.proto

package remote

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

type TestSuites struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*TestSuite `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *TestSuites) Reset() {
	*x = TestSuites{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestSuites) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestSuites) ProtoMessage() {}

func (x *TestSuites) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[0]
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
	return file_remote_loader_proto_rawDescGZIP(), []int{0}
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

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Api  string `protobuf:"bytes,2,opt,name=api,proto3" json:"api,omitempty"`
}

func (x *TestSuite) Reset() {
	*x = TestSuite{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestSuite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestSuite) ProtoMessage() {}

func (x *TestSuite) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[1]
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
	return file_remote_loader_proto_rawDescGZIP(), []int{1}
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

type TestCases struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*TestCase `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *TestCases) Reset() {
	*x = TestCases{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestCases) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCases) ProtoMessage() {}

func (x *TestCases) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCases.ProtoReflect.Descriptor instead.
func (*TestCases) Descriptor() ([]byte, []int) {
	return file_remote_loader_proto_rawDescGZIP(), []int{2}
}

func (x *TestCases) GetData() []*TestCase {
	if x != nil {
		return x.Data
	}
	return nil
}

type TestCase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SuiteName string    `protobuf:"bytes,2,opt,name=suiteName,proto3" json:"suiteName,omitempty"`
	Request   *Request  `protobuf:"bytes,3,opt,name=request,proto3" json:"request,omitempty"`
	Response  *Response `protobuf:"bytes,4,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *TestCase) Reset() {
	*x = TestCase{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCase) ProtoMessage() {}

func (x *TestCase) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCase.ProtoReflect.Descriptor instead.
func (*TestCase) Descriptor() ([]byte, []int) {
	return file_remote_loader_proto_rawDescGZIP(), []int{3}
}

func (x *TestCase) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TestCase) GetSuiteName() string {
	if x != nil {
		return x.SuiteName
	}
	return ""
}

func (x *TestCase) GetRequest() *Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *TestCase) GetResponse() *Response {
	if x != nil {
		return x.Response
	}
	return nil
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Api    string  `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
	Method string  `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Header []*Pair `protobuf:"bytes,3,rep,name=header,proto3" json:"header,omitempty"`
	Query  []*Pair `protobuf:"bytes,4,rep,name=query,proto3" json:"query,omitempty"`
	Form   []*Pair `protobuf:"bytes,5,rep,name=form,proto3" json:"form,omitempty"`
	Body   string  `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_remote_loader_proto_rawDescGZIP(), []int{4}
}

func (x *Request) GetApi() string {
	if x != nil {
		return x.Api
	}
	return ""
}

func (x *Request) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Request) GetHeader() []*Pair {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Request) GetQuery() []*Pair {
	if x != nil {
		return x.Query
	}
	return nil
}

func (x *Request) GetForm() []*Pair {
	if x != nil {
		return x.Form
	}
	return nil
}

func (x *Request) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode       int32    `protobuf:"varint,1,opt,name=statusCode,proto3" json:"statusCode,omitempty"`
	Body             string   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Header           []*Pair  `protobuf:"bytes,3,rep,name=header,proto3" json:"header,omitempty"`
	BodyFieldsExpect []*Pair  `protobuf:"bytes,4,rep,name=bodyFieldsExpect,proto3" json:"bodyFieldsExpect,omitempty"`
	Verify           []string `protobuf:"bytes,5,rep,name=verify,proto3" json:"verify,omitempty"`
	Schema           string   `protobuf:"bytes,6,opt,name=schema,proto3" json:"schema,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_remote_loader_proto_rawDescGZIP(), []int{5}
}

func (x *Response) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *Response) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Response) GetHeader() []*Pair {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Response) GetBodyFieldsExpect() []*Pair {
	if x != nil {
		return x.BodyFieldsExpect
	}
	return nil
}

func (x *Response) GetVerify() []string {
	if x != nil {
		return x.Verify
	}
	return nil
}

func (x *Response) GetSchema() string {
	if x != nil {
		return x.Schema
	}
	return ""
}

type Pair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Pair) Reset() {
	*x = Pair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pair) ProtoMessage() {}

func (x *Pair) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pair.ProtoReflect.Descriptor instead.
func (*Pair) Descriptor() ([]byte, []int) {
	return file_remote_loader_proto_rawDescGZIP(), []int{6}
}

func (x *Pair) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Pair) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_loader_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_remote_loader_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_remote_loader_proto_rawDescGZIP(), []int{7}
}

var File_remote_loader_proto protoreflect.FileDescriptor

var file_remote_loader_proto_rawDesc = []byte{
	0x0a, 0x13, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2f, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x22, 0x33, 0x0a,
	0x0a, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x31, 0x0a, 0x09, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x70, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x61, 0x70, 0x69, 0x22, 0x31, 0x0a, 0x09, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73,
	0x65, 0x73, 0x12, 0x24, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61,
	0x73, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x95, 0x01, 0x0a, 0x08, 0x54, 0x65, 0x73,
	0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x75, 0x69,
	0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75,
	0x69, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2c, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0xb3, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x61, 0x70, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x69, 0x12, 0x16,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x24, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e,
	0x50, 0x61, 0x69, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x50, 0x61, 0x69, 0x72, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x12, 0x20, 0x0a, 0x04, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x50, 0x61, 0x69, 0x72, 0x52, 0x04, 0x66, 0x6f,
	0x72, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0xce, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x24, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x2e, 0x50, 0x61, 0x69, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x38, 0x0a,
	0x10, 0x62, 0x6f, 0x64, 0x79, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x78, 0x70, 0x65, 0x63,
	0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x2e, 0x50, 0x61, 0x69, 0x72, 0x52, 0x10, 0x62, 0x6f, 0x64, 0x79, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x45, 0x78, 0x70, 0x65, 0x63, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x22, 0x2e, 0x0a, 0x04, 0x50, 0x61, 0x69, 0x72, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x32, 0xaf, 0x04, 0x0a, 0x06, 0x4c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x12, 0x34, 0x0a, 0x0d, 0x4c,
	0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x0d, 0x2e, 0x72,
	0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x73, 0x22,
	0x00, 0x12, 0x35, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x53,
	0x75, 0x69, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65,
	0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x0d, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
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
	0x65, 0x1a, 0x0d, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x37, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61,
	0x73, 0x65, 0x73, 0x12, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x1a, 0x11, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x73, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x0e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x2e,
	0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x1a,
	0x0d, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x33, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12,
	0x10, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73,
	0x65, 0x1a, 0x10, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43,
	0x61, 0x73, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x1a, 0x10, 0x2e, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x22, 0x00, 0x12, 0x33, 0x0a,
	0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12,
	0x10, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73,
	0x65, 0x1a, 0x0d, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x73, 0x75, 0x72, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2d,
	0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x69, 0x6e, 0x67, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_remote_loader_proto_rawDescOnce sync.Once
	file_remote_loader_proto_rawDescData = file_remote_loader_proto_rawDesc
)

func file_remote_loader_proto_rawDescGZIP() []byte {
	file_remote_loader_proto_rawDescOnce.Do(func() {
		file_remote_loader_proto_rawDescData = protoimpl.X.CompressGZIP(file_remote_loader_proto_rawDescData)
	})
	return file_remote_loader_proto_rawDescData
}

var file_remote_loader_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_remote_loader_proto_goTypes = []interface{}{
	(*TestSuites)(nil), // 0: remote.TestSuites
	(*TestSuite)(nil),  // 1: remote.TestSuite
	(*TestCases)(nil),  // 2: remote.TestCases
	(*TestCase)(nil),   // 3: remote.TestCase
	(*Request)(nil),    // 4: remote.Request
	(*Response)(nil),   // 5: remote.Response
	(*Pair)(nil),       // 6: remote.Pair
	(*Empty)(nil),      // 7: remote.Empty
}
var file_remote_loader_proto_depIdxs = []int32{
	1,  // 0: remote.TestSuites.data:type_name -> remote.TestSuite
	3,  // 1: remote.TestCases.data:type_name -> remote.TestCase
	4,  // 2: remote.TestCase.request:type_name -> remote.Request
	5,  // 3: remote.TestCase.response:type_name -> remote.Response
	6,  // 4: remote.Request.header:type_name -> remote.Pair
	6,  // 5: remote.Request.query:type_name -> remote.Pair
	6,  // 6: remote.Request.form:type_name -> remote.Pair
	6,  // 7: remote.Response.header:type_name -> remote.Pair
	6,  // 8: remote.Response.bodyFieldsExpect:type_name -> remote.Pair
	7,  // 9: remote.Loader.ListTestSuite:input_type -> remote.Empty
	1,  // 10: remote.Loader.CreateTestSuite:input_type -> remote.TestSuite
	1,  // 11: remote.Loader.GetTestSuite:input_type -> remote.TestSuite
	1,  // 12: remote.Loader.UpdateTestSuite:input_type -> remote.TestSuite
	1,  // 13: remote.Loader.DeleteTestSuite:input_type -> remote.TestSuite
	1,  // 14: remote.Loader.ListTestCases:input_type -> remote.TestSuite
	3,  // 15: remote.Loader.CreateTestCase:input_type -> remote.TestCase
	3,  // 16: remote.Loader.GetTestCase:input_type -> remote.TestCase
	3,  // 17: remote.Loader.UpdateTestCase:input_type -> remote.TestCase
	3,  // 18: remote.Loader.DeleteTestCase:input_type -> remote.TestCase
	0,  // 19: remote.Loader.ListTestSuite:output_type -> remote.TestSuites
	7,  // 20: remote.Loader.CreateTestSuite:output_type -> remote.Empty
	1,  // 21: remote.Loader.GetTestSuite:output_type -> remote.TestSuite
	1,  // 22: remote.Loader.UpdateTestSuite:output_type -> remote.TestSuite
	7,  // 23: remote.Loader.DeleteTestSuite:output_type -> remote.Empty
	2,  // 24: remote.Loader.ListTestCases:output_type -> remote.TestCases
	7,  // 25: remote.Loader.CreateTestCase:output_type -> remote.Empty
	3,  // 26: remote.Loader.GetTestCase:output_type -> remote.TestCase
	3,  // 27: remote.Loader.UpdateTestCase:output_type -> remote.TestCase
	7,  // 28: remote.Loader.DeleteTestCase:output_type -> remote.Empty
	19, // [19:29] is the sub-list for method output_type
	9,  // [9:19] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_remote_loader_proto_init() }
func file_remote_loader_proto_init() {
	if File_remote_loader_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_remote_loader_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_remote_loader_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_remote_loader_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestCases); i {
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
		file_remote_loader_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestCase); i {
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
		file_remote_loader_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_remote_loader_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_remote_loader_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pair); i {
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
		file_remote_loader_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_remote_loader_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_remote_loader_proto_goTypes,
		DependencyIndexes: file_remote_loader_proto_depIdxs,
		MessageInfos:      file_remote_loader_proto_msgTypes,
	}.Build()
	File_remote_loader_proto = out.File
	file_remote_loader_proto_rawDesc = nil
	file_remote_loader_proto_goTypes = nil
	file_remote_loader_proto_depIdxs = nil
}
