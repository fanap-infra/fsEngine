// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: fileIndex.proto

package fileIndex

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/descriptorpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Table struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files       map[uint32]*File `protobuf:"bytes,1,rep,name=Files,proto3" json:"Files,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	NumberFiles uint32           `protobuf:"varint,2,opt,name=NumberFiles,proto3" json:"NumberFiles,omitempty"` // number of records
}

func (x *Table) Reset() {
	*x = Table{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileIndex_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Table) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Table) ProtoMessage() {}

func (x *Table) ProtoReflect() protoreflect.Message {
	mi := &file_fileIndex_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Table.ProtoReflect.Descriptor instead.
func (*Table) Descriptor() ([]byte, []int) {
	return file_fileIndex_proto_rawDescGZIP(), []int{0}
}

func (x *Table) GetFiles() map[uint32]*File {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *Table) GetNumberFiles() uint32 {
	if x != nil {
		return x.NumberFiles
	}
	return 0
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                 // id
	FirstBlock  uint32                 `protobuf:"varint,2,opt,name=firstBlock,proto3" json:"firstBlock,omitempty"` // first block number
	LastBlock   uint32                 `protobuf:"varint,3,opt,name=lastBlock,proto3" json:"lastBlock,omitempty"`   // last block number
	RMapBlocks  []byte                 `protobuf:"bytes,4,opt,name=rMapBlocks,proto3" json:"rMapBlocks,omitempty"`  // roaring bitmap
	Name        string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	CreatedTime *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=createdTime,proto3" json:"createdTime,omitempty"`
	FileSize    uint32                 `protobuf:"varint,7,opt,name=fileSize,proto3" json:"fileSize,omitempty"`
	Optional    []byte                 `protobuf:"bytes,8,opt,name=optional,proto3" json:"optional,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileIndex_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_fileIndex_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_fileIndex_proto_rawDescGZIP(), []int{1}
}

func (x *File) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *File) GetFirstBlock() uint32 {
	if x != nil {
		return x.FirstBlock
	}
	return 0
}

func (x *File) GetLastBlock() uint32 {
	if x != nil {
		return x.LastBlock
	}
	return 0
}

func (x *File) GetRMapBlocks() []byte {
	if x != nil {
		return x.RMapBlocks
	}
	return nil
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetCreatedTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedTime
	}
	return nil
}

func (x *File) GetFileSize() uint32 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

func (x *File) GetOptional() []byte {
	if x != nil {
		return x.Optional
	}
	return nil
}

var File_fileIndex_proto protoreflect.FileDescriptor

var file_fileIndex_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x1a, 0x20, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xa7, 0x01, 0x0a, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x31, 0x0a, 0x05, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x49,
	0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x25,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xfe, 0x01, 0x0a, 0x04, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x4d, 0x61, 0x70, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x72, 0x4d, 0x61, 0x70, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_fileIndex_proto_rawDescOnce sync.Once
	file_fileIndex_proto_rawDescData = file_fileIndex_proto_rawDesc
)

func file_fileIndex_proto_rawDescGZIP() []byte {
	file_fileIndex_proto_rawDescOnce.Do(func() {
		file_fileIndex_proto_rawDescData = protoimpl.X.CompressGZIP(file_fileIndex_proto_rawDescData)
	})
	return file_fileIndex_proto_rawDescData
}

var (
	file_fileIndex_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
	file_fileIndex_proto_goTypes  = []interface{}{
		(*Table)(nil),                 // 0: fileIndex.Table
		(*File)(nil),                  // 1: fileIndex.File
		nil,                           // 2: fileIndex.Table.FilesEntry
		(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	}
)

var file_fileIndex_proto_depIdxs = []int32{
	2, // 0: fileIndex.Table.Files:type_name -> fileIndex.Table.FilesEntry
	3, // 1: fileIndex.File.createdTime:type_name -> google.protobuf.Timestamp
	1, // 2: fileIndex.Table.FilesEntry.value:type_name -> fileIndex.File
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_fileIndex_proto_init() }
func file_fileIndex_proto_init() {
	if File_fileIndex_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fileIndex_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Table); i {
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
		file_fileIndex_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
			RawDescriptor: file_fileIndex_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fileIndex_proto_goTypes,
		DependencyIndexes: file_fileIndex_proto_depIdxs,
		MessageInfos:      file_fileIndex_proto_msgTypes,
	}.Build()
	File_fileIndex_proto = out.File
	file_fileIndex_proto_rawDesc = nil
	file_fileIndex_proto_goTypes = nil
	file_fileIndex_proto_depIdxs = nil
}