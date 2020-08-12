// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.11.4
// source: generated.proto

package protobuf

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type RequestInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method      string `protobuf:"bytes,1,opt,name=Method,proto3" json:"Method,omitempty"`
	Path        string `protobuf:"bytes,2,opt,name=Path,proto3" json:"Path,omitempty"`
	ClusterName string `protobuf:"bytes,3,opt,name=ClusterName,proto3" json:"ClusterName,omitempty"`
	Body        string `protobuf:"bytes,4,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *RequestInfo) Reset() {
	*x = RequestInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generated_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestInfo) ProtoMessage() {}

func (x *RequestInfo) ProtoReflect() protoreflect.Message {
	mi := &file_generated_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestInfo.ProtoReflect.Descriptor instead.
func (*RequestInfo) Descriptor() ([]byte, []int) {
	return file_generated_proto_rawDescGZIP(), []int{0}
}

func (x *RequestInfo) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *RequestInfo) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *RequestInfo) GetClusterName() string {
	if x != nil {
		return x.ClusterName
	}
	return ""
}

func (x *RequestInfo) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type ResponseInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  string `protobuf:"bytes,1,opt,name=Status,proto3" json:"Status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *ResponseInfo) Reset() {
	*x = ResponseInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_generated_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseInfo) ProtoMessage() {}

func (x *ResponseInfo) ProtoReflect() protoreflect.Message {
	mi := &file_generated_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseInfo.ProtoReflect.Descriptor instead.
func (*ResponseInfo) Descriptor() ([]byte, []int) {
	return file_generated_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseInfo) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ResponseInfo) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_generated_proto protoreflect.FileDescriptor

var file_generated_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x6f, 0x0a, 0x0b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x50, 0x61, 0x74, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x40, 0x0a, 0x0c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x5b,
	0x0a, 0x10, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x50, 0x49, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x12, 0x47, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x4f, 0x70, 0x65, 0x6e, 0x4d, 0x43,
	0x50, 0x41, 0x50, 0x49, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x2e,
	0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_generated_proto_rawDescOnce sync.Once
	file_generated_proto_rawDescData = file_generated_proto_rawDesc
)

func file_generated_proto_rawDescGZIP() []byte {
	file_generated_proto_rawDescOnce.Do(func() {
		file_generated_proto_rawDescData = protoimpl.X.CompressGZIP(file_generated_proto_rawDescData)
	})
	return file_generated_proto_rawDescData
}

var file_generated_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_generated_proto_goTypes = []interface{}{
	(*RequestInfo)(nil),  // 0: protobuf.RequestInfo
	(*ResponseInfo)(nil), // 1: protobuf.ResponseInfo
}
var file_generated_proto_depIdxs = []int32{
	0, // 0: protobuf.RequestAPIServer.SendOpenMCPAPIServer:input_type -> protobuf.RequestInfo
	1, // 1: protobuf.RequestAPIServer.SendOpenMCPAPIServer:output_type -> protobuf.ResponseInfo
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_generated_proto_init() }
func file_generated_proto_init() {
	if File_generated_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_generated_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestInfo); i {
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
		file_generated_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseInfo); i {
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
			RawDescriptor: file_generated_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_generated_proto_goTypes,
		DependencyIndexes: file_generated_proto_depIdxs,
		MessageInfos:      file_generated_proto_msgTypes,
	}.Build()
	File_generated_proto = out.File
	file_generated_proto_rawDesc = nil
	file_generated_proto_goTypes = nil
	file_generated_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RequestAPIServerClient is the client API for RequestAPIServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RequestAPIServerClient interface {
	SendOpenMCPAPIServer(ctx context.Context, in *RequestInfo, opts ...grpc.CallOption) (*ResponseInfo, error)
}

type requestAPIServerClient struct {
	cc grpc.ClientConnInterface
}

func NewRequestAPIServerClient(cc grpc.ClientConnInterface) RequestAPIServerClient {
	return &requestAPIServerClient{cc}
}

func (c *requestAPIServerClient) SendOpenMCPAPIServer(ctx context.Context, in *RequestInfo, opts ...grpc.CallOption) (*ResponseInfo, error) {
	out := new(ResponseInfo)
	err := c.cc.Invoke(ctx, "/protobuf.RequestAPIServer/SendOpenMCPAPIServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RequestAPIServerServer is the server API for RequestAPIServer service.
type RequestAPIServerServer interface {
	SendOpenMCPAPIServer(context.Context, *RequestInfo) (*ResponseInfo, error)
}

// UnimplementedRequestAPIServerServer can be embedded to have forward compatible implementations.
type UnimplementedRequestAPIServerServer struct {
}

func (*UnimplementedRequestAPIServerServer) SendOpenMCPAPIServer(context.Context, *RequestInfo) (*ResponseInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendOpenMCPAPIServer not implemented")
}

func RegisterRequestAPIServerServer(s *grpc.Server, srv RequestAPIServerServer) {
	s.RegisterService(&_RequestAPIServer_serviceDesc, srv)
}

func _RequestAPIServer_SendOpenMCPAPIServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestAPIServerServer).SendOpenMCPAPIServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.RequestAPIServer/SendOpenMCPAPIServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestAPIServerServer).SendOpenMCPAPIServer(ctx, req.(*RequestInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _RequestAPIServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.RequestAPIServer",
	HandlerType: (*RequestAPIServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendOpenMCPAPIServer",
			Handler:    _RequestAPIServer_SendOpenMCPAPIServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "generated.proto",
}