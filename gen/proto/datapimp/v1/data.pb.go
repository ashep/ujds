// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: proto/datapimp/v1/data.proto

package v1

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

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type    string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Version uint64 `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	Time    uint64 `protobuf:"varint,4,opt,name=time,proto3" json:"time,omitempty"`
	Data    string `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_datapimp_v1_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_proto_datapimp_v1_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_proto_datapimp_v1_data_proto_rawDescGZIP(), []int{0}
}

func (x *Item) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Item) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Item) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Item) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Item) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type PushItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Data string `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *PushItemRequest) Reset() {
	*x = PushItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_datapimp_v1_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushItemRequest) ProtoMessage() {}

func (x *PushItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_datapimp_v1_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushItemRequest.ProtoReflect.Descriptor instead.
func (*PushItemRequest) Descriptor() ([]byte, []int) {
	return file_proto_datapimp_v1_data_proto_rawDescGZIP(), []int{1}
}

func (x *PushItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PushItemRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *PushItemRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type PushItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item *Item `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
}

func (x *PushItemResponse) Reset() {
	*x = PushItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_datapimp_v1_data_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushItemResponse) ProtoMessage() {}

func (x *PushItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_datapimp_v1_data_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushItemResponse.ProtoReflect.Descriptor instead.
func (*PushItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_datapimp_v1_data_proto_rawDescGZIP(), []int{2}
}

func (x *PushItemResponse) GetItem() *Item {
	if x != nil {
		return x.Item
	}
	return nil
}

type GetItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetItemRequest) Reset() {
	*x = GetItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_datapimp_v1_data_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemRequest) ProtoMessage() {}

func (x *GetItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_datapimp_v1_data_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemRequest.ProtoReflect.Descriptor instead.
func (*GetItemRequest) Descriptor() ([]byte, []int) {
	return file_proto_datapimp_v1_data_proto_rawDescGZIP(), []int{3}
}

func (x *GetItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Item *Item `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
}

func (x *GetItemResponse) Reset() {
	*x = GetItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_datapimp_v1_data_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemResponse) ProtoMessage() {}

func (x *GetItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_datapimp_v1_data_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemResponse.ProtoReflect.Descriptor instead.
func (*GetItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_datapimp_v1_data_proto_rawDescGZIP(), []int{4}
}

func (x *GetItemResponse) GetItem() *Item {
	if x != nil {
		return x.Item
	}
	return nil
}

var File_proto_datapimp_v1_data_proto protoreflect.FileDescriptor

var file_proto_datapimp_v1_data_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70,
	0x2f, 0x76, 0x31, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70, 0x2e, 0x76, 0x31, 0x22, 0x6c, 0x0a, 0x04, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x49, 0x0a, 0x0f, 0x50, 0x75, 0x73,
	0x68, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x39, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d,
	0x70, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x22,
	0x20, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x38, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x32, 0xa0, 0x01, 0x0a, 0x0b,
	0x44, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x08, 0x50,
	0x75, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x1c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69,
	0x6d, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x31,
	0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x73, 0x68,
	0x65, 0x70, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x69, 0x6d, 0x70, 0x2f, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_datapimp_v1_data_proto_rawDescOnce sync.Once
	file_proto_datapimp_v1_data_proto_rawDescData = file_proto_datapimp_v1_data_proto_rawDesc
)

func file_proto_datapimp_v1_data_proto_rawDescGZIP() []byte {
	file_proto_datapimp_v1_data_proto_rawDescOnce.Do(func() {
		file_proto_datapimp_v1_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_datapimp_v1_data_proto_rawDescData)
	})
	return file_proto_datapimp_v1_data_proto_rawDescData
}

var file_proto_datapimp_v1_data_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_datapimp_v1_data_proto_goTypes = []interface{}{
	(*Item)(nil),             // 0: datapimp.v1.Item
	(*PushItemRequest)(nil),  // 1: datapimp.v1.PushItemRequest
	(*PushItemResponse)(nil), // 2: datapimp.v1.PushItemResponse
	(*GetItemRequest)(nil),   // 3: datapimp.v1.GetItemRequest
	(*GetItemResponse)(nil),  // 4: datapimp.v1.GetItemResponse
}
var file_proto_datapimp_v1_data_proto_depIdxs = []int32{
	0, // 0: datapimp.v1.PushItemResponse.item:type_name -> datapimp.v1.Item
	0, // 1: datapimp.v1.GetItemResponse.item:type_name -> datapimp.v1.Item
	1, // 2: datapimp.v1.DataService.PushItem:input_type -> datapimp.v1.PushItemRequest
	3, // 3: datapimp.v1.DataService.GetItem:input_type -> datapimp.v1.GetItemRequest
	2, // 4: datapimp.v1.DataService.PushItem:output_type -> datapimp.v1.PushItemResponse
	4, // 5: datapimp.v1.DataService.GetItem:output_type -> datapimp.v1.GetItemResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_datapimp_v1_data_proto_init() }
func file_proto_datapimp_v1_data_proto_init() {
	if File_proto_datapimp_v1_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_datapimp_v1_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Item); i {
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
		file_proto_datapimp_v1_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushItemRequest); i {
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
		file_proto_datapimp_v1_data_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushItemResponse); i {
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
		file_proto_datapimp_v1_data_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetItemRequest); i {
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
		file_proto_datapimp_v1_data_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetItemResponse); i {
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
			RawDescriptor: file_proto_datapimp_v1_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_datapimp_v1_data_proto_goTypes,
		DependencyIndexes: file_proto_datapimp_v1_data_proto_depIdxs,
		MessageInfos:      file_proto_datapimp_v1_data_proto_msgTypes,
	}.Build()
	File_proto_datapimp_v1_data_proto = out.File
	file_proto_datapimp_v1_data_proto_rawDesc = nil
	file_proto_datapimp_v1_data_proto_goTypes = nil
	file_proto_datapimp_v1_data_proto_depIdxs = nil
}
