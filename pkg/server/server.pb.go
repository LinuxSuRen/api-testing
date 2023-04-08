// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/server/server.proto

package server

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TestTask struct {
	Data                 string            `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Kind                 string            `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	CaseName             string            `protobuf:"bytes,3,opt,name=caseName,proto3" json:"caseName,omitempty"`
	Level                string            `protobuf:"bytes,4,opt,name=level,proto3" json:"level,omitempty"`
	Env                  map[string]string `protobuf:"bytes,5,rep,name=env,proto3" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TestTask) Reset()         { *m = TestTask{} }
func (m *TestTask) String() string { return proto.CompactTextString(m) }
func (*TestTask) ProtoMessage()    {}
func (*TestTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_36fb7b77b8f76c98, []int{0}
}

func (m *TestTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestTask.Unmarshal(m, b)
}
func (m *TestTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestTask.Marshal(b, m, deterministic)
}
func (m *TestTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestTask.Merge(m, src)
}
func (m *TestTask) XXX_Size() int {
	return xxx_messageInfo_TestTask.Size(m)
}
func (m *TestTask) XXX_DiscardUnknown() {
	xxx_messageInfo_TestTask.DiscardUnknown(m)
}

var xxx_messageInfo_TestTask proto.InternalMessageInfo

func (m *TestTask) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *TestTask) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *TestTask) GetCaseName() string {
	if m != nil {
		return m.CaseName
	}
	return ""
}

func (m *TestTask) GetLevel() string {
	if m != nil {
		return m.Level
	}
	return ""
}

func (m *TestTask) GetEnv() map[string]string {
	if m != nil {
		return m.Env
	}
	return nil
}

type HelloReply struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloReply) Reset()         { *m = HelloReply{} }
func (m *HelloReply) String() string { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()    {}
func (*HelloReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_36fb7b77b8f76c98, []int{1}
}

func (m *HelloReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloReply.Unmarshal(m, b)
}
func (m *HelloReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloReply.Marshal(b, m, deterministic)
}
func (m *HelloReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReply.Merge(m, src)
}
func (m *HelloReply) XXX_Size() int {
	return xxx_messageInfo_HelloReply.Size(m)
}
func (m *HelloReply) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReply.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReply proto.InternalMessageInfo

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_36fb7b77b8f76c98, []int{2}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TestTask)(nil), "server.TestTask")
	proto.RegisterMapType((map[string]string)(nil), "server.TestTask.EnvEntry")
	proto.RegisterType((*HelloReply)(nil), "server.HelloReply")
	proto.RegisterType((*Empty)(nil), "server.Empty")
}

func init() {
	proto.RegisterFile("pkg/server/server.proto", fileDescriptor_36fb7b77b8f76c98)
}

var fileDescriptor_36fb7b77b8f76c98 = []byte{
	// 301 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x4b, 0x4b, 0xfb, 0x40,
	0x14, 0xc5, 0xff, 0x69, 0xfa, 0xfa, 0x5f, 0x11, 0xca, 0x45, 0x70, 0xec, 0xaa, 0x64, 0x21, 0x05,
	0x6d, 0x82, 0x15, 0x44, 0x5c, 0x0a, 0x45, 0x57, 0x2e, 0x42, 0x71, 0xe1, 0x6e, 0xda, 0x5e, 0xe2,
	0x98, 0xc9, 0x24, 0xcc, 0x23, 0x98, 0x8f, 0xe8, 0xb7, 0x92, 0xbc, 0x2c, 0x88, 0xab, 0x39, 0xe7,
	0x37, 0xf7, 0x71, 0xe0, 0xc2, 0x79, 0x91, 0x26, 0x91, 0x21, 0x5d, 0x92, 0xee, 0x9e, 0xb0, 0xd0,
	0xb9, 0xcd, 0x71, 0xdc, 0xba, 0xe0, 0xcb, 0x83, 0xe9, 0x96, 0x8c, 0xdd, 0x72, 0x93, 0x22, 0xc2,
	0xf0, 0xc0, 0x2d, 0x67, 0xde, 0xc2, 0x5b, 0xfe, 0x8f, 0x1b, 0x5d, 0xb3, 0x54, 0xa8, 0x03, 0x1b,
	0xb4, 0xac, 0xd6, 0x38, 0x87, 0xe9, 0x9e, 0x1b, 0x7a, 0xe1, 0x19, 0x31, 0xbf, 0xe1, 0x3f, 0x1e,
	0xcf, 0x60, 0x24, 0xa9, 0x24, 0xc9, 0x86, 0xcd, 0x47, 0x6b, 0xf0, 0x0a, 0x7c, 0x52, 0x25, 0x1b,
	0x2d, 0xfc, 0xe5, 0xc9, 0xfa, 0x22, 0xec, 0xa2, 0xf4, 0x8b, 0xc3, 0x8d, 0x2a, 0x37, 0xca, 0xea,
	0x2a, 0xae, 0xab, 0xe6, 0x77, 0x30, 0xed, 0x01, 0xce, 0xc0, 0x4f, 0xa9, 0xea, 0x12, 0xd5, 0xb2,
	0x5e, 0x50, 0x72, 0xe9, 0xa8, 0x4b, 0xd4, 0x9a, 0x87, 0xc1, 0xbd, 0x17, 0x5c, 0x02, 0x3c, 0x93,
	0x94, 0x79, 0x4c, 0x85, 0xac, 0x90, 0xc1, 0x24, 0x23, 0x63, 0x78, 0x42, 0x5d, 0x77, 0x6f, 0x83,
	0x09, 0x8c, 0x36, 0x59, 0x61, 0xab, 0xf5, 0x07, 0x8c, 0x63, 0xa7, 0x14, 0x69, 0x5c, 0x81, 0x1f,
	0x3b, 0x85, 0xb3, 0xdf, 0xc9, 0xe6, 0xd8, 0x93, 0xe3, 0xe4, 0xe0, 0x1f, 0xde, 0x00, 0x3c, 0x91,
	0x7d, 0x25, 0x6d, 0x44, 0xae, 0xf0, 0xb4, 0xaf, 0x69, 0xa6, 0xfe, 0xdd, 0xf2, 0x18, 0xbe, 0x5d,
	0x27, 0xc2, 0xbe, 0xbb, 0x5d, 0xb8, 0xcf, 0xb3, 0x48, 0x0a, 0xe5, 0x3e, 0x8d, 0xd3, 0xa4, 0x22,
	0x5e, 0x88, 0x95, 0x25, 0x63, 0x85, 0x4a, 0xa2, 0xe3, 0xb5, 0x76, 0xe3, 0xe6, 0x4e, 0xb7, 0xdf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xa2, 0x3b, 0xc8, 0x00, 0xc2, 0x01, 0x00, 0x00,
}
