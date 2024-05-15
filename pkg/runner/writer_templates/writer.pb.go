// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.25.3
// source: writer_templates/writer.proto

package writer_templates

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

type ReportResultRepeated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*ReportResult `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ReportResultRepeated) Reset() {
	*x = ReportResultRepeated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_writer_templates_writer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportResultRepeated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportResultRepeated) ProtoMessage() {}

func (x *ReportResultRepeated) ProtoReflect() protoreflect.Message {
	mi := &file_writer_templates_writer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportResultRepeated.ProtoReflect.Descriptor instead.
func (*ReportResultRepeated) Descriptor() ([]byte, []int) {
	return file_writer_templates_writer_proto_rawDescGZIP(), []int{0}
}

func (x *ReportResultRepeated) GetData() []*ReportResult {
	if x != nil {
		return x.Data
	}
	return nil
}

type ReportResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name             string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	API              string `protobuf:"bytes,2,opt,name=API,proto3" json:"API,omitempty"`
	Count            int32  `protobuf:"varint,3,opt,name=Count,proto3" json:"Count,omitempty"`
	Average          int64  `protobuf:"varint,4,opt,name=Average,proto3" json:"Average,omitempty"`
	Max              int64  `protobuf:"varint,5,opt,name=Max,proto3" json:"Max,omitempty"`
	Min              int64  `protobuf:"varint,6,opt,name=Min,proto3" json:"Min,omitempty"`
	QPS              int32  `protobuf:"varint,7,opt,name=QPS,proto3" json:"QPS,omitempty"`
	Error            int32  `protobuf:"varint,8,opt,name=Error,proto3" json:"Error,omitempty"`
	LastErrorMessage string `protobuf:"bytes,9,opt,name=LastErrorMessage,proto3" json:"LastErrorMessage,omitempty"`
}

func (x *ReportResult) Reset() {
	*x = ReportResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_writer_templates_writer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportResult) ProtoMessage() {}

func (x *ReportResult) ProtoReflect() protoreflect.Message {
	mi := &file_writer_templates_writer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportResult.ProtoReflect.Descriptor instead.
func (*ReportResult) Descriptor() ([]byte, []int) {
	return file_writer_templates_writer_proto_rawDescGZIP(), []int{1}
}

func (x *ReportResult) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ReportResult) GetAPI() string {
	if x != nil {
		return x.API
	}
	return ""
}

func (x *ReportResult) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *ReportResult) GetAverage() int64 {
	if x != nil {
		return x.Average
	}
	return 0
}

func (x *ReportResult) GetMax() int64 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *ReportResult) GetMin() int64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *ReportResult) GetQPS() int32 {
	if x != nil {
		return x.QPS
	}
	return 0
}

func (x *ReportResult) GetError() int32 {
	if x != nil {
		return x.Error
	}
	return 0
}

func (x *ReportResult) GetLastErrorMessage() string {
	if x != nil {
		return x.LastErrorMessage
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
		mi := &file_writer_templates_writer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_writer_templates_writer_proto_msgTypes[2]
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
	return file_writer_templates_writer_proto_rawDescGZIP(), []int{2}
}

var File_writer_templates_writer_proto protoreflect.FileDescriptor

var file_writer_templates_writer_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x73, 0x2f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x10, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x73, 0x22, 0x4a, 0x0a, 0x14, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x32, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72,
	0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xdc, 0x01,
	0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x41, 0x50, 0x49, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x76,
	0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x41, 0x76, 0x65,
	0x72, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x61, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x4d, 0x61, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x69, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x03, 0x4d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x51, 0x50, 0x53, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x51, 0x50, 0x53, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x2a, 0x0a, 0x10, 0x4c, 0x61, 0x73, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x4c, 0x61, 0x73, 0x74,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x07, 0x0a, 0x05,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x63, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x72, 0x12, 0x53, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x26, 0x2e, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x1a, 0x17, 0x2e, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c,
	0x61, 0x74, 0x65, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x73, 0x75,
	0x72, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x77, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_writer_templates_writer_proto_rawDescOnce sync.Once
	file_writer_templates_writer_proto_rawDescData = file_writer_templates_writer_proto_rawDesc
)

func file_writer_templates_writer_proto_rawDescGZIP() []byte {
	file_writer_templates_writer_proto_rawDescOnce.Do(func() {
		file_writer_templates_writer_proto_rawDescData = protoimpl.X.CompressGZIP(file_writer_templates_writer_proto_rawDescData)
	})
	return file_writer_templates_writer_proto_rawDescData
}

var file_writer_templates_writer_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_writer_templates_writer_proto_goTypes = []interface{}{
	(*ReportResultRepeated)(nil), // 0: writer_templates.ReportResultRepeated
	(*ReportResult)(nil),         // 1: writer_templates.ReportResult
	(*Empty)(nil),                // 2: writer_templates.Empty
}
var file_writer_templates_writer_proto_depIdxs = []int32{
	1, // 0: writer_templates.ReportResultRepeated.data:type_name -> writer_templates.ReportResult
	0, // 1: writer_templates.ReportWriter.SendReportResult:input_type -> writer_templates.ReportResultRepeated
	2, // 2: writer_templates.ReportWriter.SendReportResult:output_type -> writer_templates.Empty
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_writer_templates_writer_proto_init() }
func file_writer_templates_writer_proto_init() {
	if File_writer_templates_writer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_writer_templates_writer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportResultRepeated); i {
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
		file_writer_templates_writer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportResult); i {
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
		file_writer_templates_writer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_writer_templates_writer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_writer_templates_writer_proto_goTypes,
		DependencyIndexes: file_writer_templates_writer_proto_depIdxs,
		MessageInfos:      file_writer_templates_writer_proto_msgTypes,
	}.Build()
	File_writer_templates_writer_proto = out.File
	file_writer_templates_writer_proto_rawDesc = nil
	file_writer_templates_writer_proto_goTypes = nil
	file_writer_templates_writer_proto_depIdxs = nil
}
