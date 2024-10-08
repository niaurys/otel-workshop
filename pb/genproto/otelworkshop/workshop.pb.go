// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.2
// source: workshop.proto

package otelworkshop

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workshop_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_workshop_proto_msgTypes[0]
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
	return file_workshop_proto_rawDescGZIP(), []int{0}
}

type Product struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Color    string `protobuf:"bytes,2,opt,name=color,proto3" json:"color,omitempty"`
	Quantity int64  `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *Product) Reset() {
	*x = Product{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workshop_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_workshop_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_workshop_proto_rawDescGZIP(), []int{1}
}

func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Product) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *Product) GetQuantity() int64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type ListProductsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Products []*Product `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
}

func (x *ListProductsResponse) Reset() {
	*x = ListProductsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workshop_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListProductsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProductsResponse) ProtoMessage() {}

func (x *ListProductsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_workshop_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProductsResponse.ProtoReflect.Descriptor instead.
func (*ListProductsResponse) Descriptor() ([]byte, []int) {
	return file_workshop_proto_rawDescGZIP(), []int{2}
}

func (x *ListProductsResponse) GetProducts() []*Product {
	if x != nil {
		return x.Products
	}
	return nil
}

type BuyProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Surname string   `protobuf:"bytes,2,opt,name=surname,proto3" json:"surname,omitempty"`
	Product *Product `protobuf:"bytes,3,opt,name=product,proto3" json:"product,omitempty"`
}

func (x *BuyProductRequest) Reset() {
	*x = BuyProductRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workshop_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuyProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuyProductRequest) ProtoMessage() {}

func (x *BuyProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workshop_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuyProductRequest.ProtoReflect.Descriptor instead.
func (*BuyProductRequest) Descriptor() ([]byte, []int) {
	return file_workshop_proto_rawDescGZIP(), []int{3}
}

func (x *BuyProductRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BuyProductRequest) GetSurname() string {
	if x != nil {
		return x.Surname
	}
	return ""
}

func (x *BuyProductRequest) GetProduct() *Product {
	if x != nil {
		return x.Product
	}
	return nil
}

var File_workshop_proto protoreflect.FileDescriptor

var file_workshop_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0c, 0x6f, 0x74, 0x65, 0x6c, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x6f, 0x70, 0x22, 0x07,
	0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x4f, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x1a, 0x0a, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x49, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x31, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x6f,
	0x70, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x73, 0x22, 0x72, 0x0a, 0x11, 0x42, 0x75, 0x79, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x77, 0x6f,
	0x72, 0x6b, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x07,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x32, 0xa0, 0x01, 0x0a, 0x0b, 0x53, 0x68, 0x6f, 0x70,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x12, 0x13, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x77, 0x6f,
	0x72, 0x6b, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x22, 0x2e, 0x6f,
	0x74, 0x65, 0x6c, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x46, 0x0a, 0x0a, 0x42, 0x75, 0x79, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x12, 0x1f, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x6f, 0x70, 0x2e,
	0x42, 0x75, 0x79, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x6f, 0x74, 0x65, 0x6c, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x6f, 0x70,
	0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x22, 0x00, 0x42, 0x17, 0x5a, 0x15, 0x67, 0x65,
	0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x74, 0x65, 0x6c, 0x77, 0x6f, 0x72, 0x6b, 0x73,
	0x68, 0x6f, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workshop_proto_rawDescOnce sync.Once
	file_workshop_proto_rawDescData = file_workshop_proto_rawDesc
)

func file_workshop_proto_rawDescGZIP() []byte {
	file_workshop_proto_rawDescOnce.Do(func() {
		file_workshop_proto_rawDescData = protoimpl.X.CompressGZIP(file_workshop_proto_rawDescData)
	})
	return file_workshop_proto_rawDescData
}

var file_workshop_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_workshop_proto_goTypes = []any{
	(*Empty)(nil),                // 0: otelworkshop.Empty
	(*Product)(nil),              // 1: otelworkshop.Product
	(*ListProductsResponse)(nil), // 2: otelworkshop.ListProductsResponse
	(*BuyProductRequest)(nil),    // 3: otelworkshop.BuyProductRequest
}
var file_workshop_proto_depIdxs = []int32{
	1, // 0: otelworkshop.ListProductsResponse.products:type_name -> otelworkshop.Product
	1, // 1: otelworkshop.BuyProductRequest.product:type_name -> otelworkshop.Product
	0, // 2: otelworkshop.ShopService.ListProducts:input_type -> otelworkshop.Empty
	3, // 3: otelworkshop.ShopService.BuyProduct:input_type -> otelworkshop.BuyProductRequest
	2, // 4: otelworkshop.ShopService.ListProducts:output_type -> otelworkshop.ListProductsResponse
	1, // 5: otelworkshop.ShopService.BuyProduct:output_type -> otelworkshop.Product
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_workshop_proto_init() }
func file_workshop_proto_init() {
	if File_workshop_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_workshop_proto_msgTypes[0].Exporter = func(v any, i int) any {
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
		file_workshop_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Product); i {
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
		file_workshop_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ListProductsResponse); i {
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
		file_workshop_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*BuyProductRequest); i {
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
			RawDescriptor: file_workshop_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_workshop_proto_goTypes,
		DependencyIndexes: file_workshop_proto_depIdxs,
		MessageInfos:      file_workshop_proto_msgTypes,
	}.Build()
	File_workshop_proto = out.File
	file_workshop_proto_rawDesc = nil
	file_workshop_proto_goTypes = nil
	file_workshop_proto_depIdxs = nil
}
