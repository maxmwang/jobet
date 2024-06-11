// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: internal/proto/jobet.proto

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

type ProbeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name of the company to probe
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// dry run. do not add to the database
	Dry bool `protobuf:"varint,2,opt,name=dry,proto3" json:"dry,omitempty"`
	// alias to use for the company
	Alias *string `protobuf:"bytes,3,opt,name=alias,proto3,oneof" json:"alias,omitempty"`
	// priority of the company. lower value is higher priority
	Priority *uint32 `protobuf:"varint,4,opt,name=priority,proto3,oneof" json:"priority,omitempty"`
}

func (x *ProbeRequest) Reset() {
	*x = ProbeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_jobet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeRequest) ProtoMessage() {}

func (x *ProbeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_jobet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeRequest.ProtoReflect.Descriptor instead.
func (*ProbeRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_jobet_proto_rawDescGZIP(), []int{0}
}

func (x *ProbeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProbeRequest) GetDry() bool {
	if x != nil {
		return x.Dry
	}
	return false
}

func (x *ProbeRequest) GetAlias() string {
	if x != nil && x.Alias != nil {
		return *x.Alias
	}
	return ""
}

func (x *ProbeRequest) GetPriority() uint32 {
	if x != nil && x.Priority != nil {
		return *x.Priority
	}
	return 0
}

type ProbeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results []*ProbeReply_ProbeSiteResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *ProbeReply) Reset() {
	*x = ProbeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_jobet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeReply) ProtoMessage() {}

func (x *ProbeReply) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_jobet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeReply.ProtoReflect.Descriptor instead.
func (*ProbeReply) Descriptor() ([]byte, []int) {
	return file_internal_proto_jobet_proto_rawDescGZIP(), []int{1}
}

func (x *ProbeReply) GetResults() []*ProbeReply_ProbeSiteResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type ProbeReply_ProbeSiteResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// listing site
	Site string `protobuf:"bytes,1,opt,name=site,proto3" json:"site,omitempty"`
	// added to the database. always false if dry run
	Added bool `protobuf:"varint,2,opt,name=added,proto3" json:"added,omitempty"`
	// already exists in the database
	Exists bool `protobuf:"varint,3,opt,name=exists,proto3" json:"exists,omitempty"`
	// number of listings found
	Count uint32 `protobuf:"varint,4,opt,name=count,proto3" json:"count,omitempty"`
	// number of target listings found
	Target uint32 `protobuf:"varint,5,opt,name=target,proto3" json:"target,omitempty"`
}

func (x *ProbeReply_ProbeSiteResult) Reset() {
	*x = ProbeReply_ProbeSiteResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_jobet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeReply_ProbeSiteResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeReply_ProbeSiteResult) ProtoMessage() {}

func (x *ProbeReply_ProbeSiteResult) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_jobet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeReply_ProbeSiteResult.ProtoReflect.Descriptor instead.
func (*ProbeReply_ProbeSiteResult) Descriptor() ([]byte, []int) {
	return file_internal_proto_jobet_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ProbeReply_ProbeSiteResult) GetSite() string {
	if x != nil {
		return x.Site
	}
	return ""
}

func (x *ProbeReply_ProbeSiteResult) GetAdded() bool {
	if x != nil {
		return x.Added
	}
	return false
}

func (x *ProbeReply_ProbeSiteResult) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

func (x *ProbeReply_ProbeSiteResult) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *ProbeReply_ProbeSiteResult) GetTarget() uint32 {
	if x != nil {
		return x.Target
	}
	return 0
}

var File_internal_proto_jobet_proto protoreflect.FileDescriptor

var file_internal_proto_jobet_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x6a, 0x6f, 0x62, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x87, 0x01, 0x0a,
	0x0c, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03,
	0x64, 0x72, 0x79, 0x12, 0x19, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x88, 0x01, 0x01, 0x12, 0x1f,
	0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x48, 0x01, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x22, 0xc7, 0x01, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x62, 0x65,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x35, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x53, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x1a, 0x81, 0x01, 0x0a,
	0x0f, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x53, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x73, 0x69, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x64, 0x64, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x64, 0x64, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78,
	0x69, 0x73, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x78, 0x69, 0x73,
	0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x32, 0x2e, 0x0a, 0x05, 0x4a, 0x6f, 0x62, 0x65, 0x74, 0x12, 0x25, 0x0a, 0x05, 0x50, 0x72, 0x6f,
	0x62, 0x65, 0x12, 0x0d, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0b, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x42, 0x12, 0x5a, 0x10, 0x2e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_jobet_proto_rawDescOnce sync.Once
	file_internal_proto_jobet_proto_rawDescData = file_internal_proto_jobet_proto_rawDesc
)

func file_internal_proto_jobet_proto_rawDescGZIP() []byte {
	file_internal_proto_jobet_proto_rawDescOnce.Do(func() {
		file_internal_proto_jobet_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_jobet_proto_rawDescData)
	})
	return file_internal_proto_jobet_proto_rawDescData
}

var file_internal_proto_jobet_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_proto_jobet_proto_goTypes = []interface{}{
	(*ProbeRequest)(nil),               // 0: ProbeRequest
	(*ProbeReply)(nil),                 // 1: ProbeReply
	(*ProbeReply_ProbeSiteResult)(nil), // 2: ProbeReply.ProbeSiteResult
}
var file_internal_proto_jobet_proto_depIdxs = []int32{
	2, // 0: ProbeReply.results:type_name -> ProbeReply.ProbeSiteResult
	0, // 1: Jobet.Probe:input_type -> ProbeRequest
	1, // 2: Jobet.Probe:output_type -> ProbeReply
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_proto_jobet_proto_init() }
func file_internal_proto_jobet_proto_init() {
	if File_internal_proto_jobet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_jobet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeRequest); i {
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
		file_internal_proto_jobet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeReply); i {
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
		file_internal_proto_jobet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeReply_ProbeSiteResult); i {
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
	file_internal_proto_jobet_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_proto_jobet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_jobet_proto_goTypes,
		DependencyIndexes: file_internal_proto_jobet_proto_depIdxs,
		MessageInfos:      file_internal_proto_jobet_proto_msgTypes,
	}.Build()
	File_internal_proto_jobet_proto = out.File
	file_internal_proto_jobet_proto_rawDesc = nil
	file_internal_proto_jobet_proto_goTypes = nil
	file_internal_proto_jobet_proto_depIdxs = nil
}
