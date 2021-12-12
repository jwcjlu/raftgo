// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

// 指定等会文件生成出来的package

package service

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// 定义request
type VoteRequest struct {
	Term                 int32    `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	NodeId               string   `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteRequest) Reset()         { *m = VoteRequest{} }
func (m *VoteRequest) String() string { return proto.CompactTextString(m) }
func (*VoteRequest) ProtoMessage()    {}
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *VoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteRequest.Unmarshal(m, b)
}
func (m *VoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteRequest.Marshal(b, m, deterministic)
}
func (m *VoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteRequest.Merge(m, src)
}
func (m *VoteRequest) XXX_Size() int {
	return xxx_messageInfo_VoteRequest.Size(m)
}
func (m *VoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VoteRequest proto.InternalMessageInfo

func (m *VoteRequest) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *VoteRequest) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

// 定义response
type Response struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Response) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type VoteResponse struct {
	Rsp                  *Response `protobuf:"bytes,1,opt,name=rsp,proto3" json:"rsp,omitempty"`
	Result               bool      `protobuf:"varint,2,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *VoteResponse) Reset()         { *m = VoteResponse{} }
func (m *VoteResponse) String() string { return proto.CompactTextString(m) }
func (*VoteResponse) ProtoMessage()    {}
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{2}
}

func (m *VoteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteResponse.Unmarshal(m, b)
}
func (m *VoteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteResponse.Marshal(b, m, deterministic)
}
func (m *VoteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteResponse.Merge(m, src)
}
func (m *VoteResponse) XXX_Size() int {
	return xxx_messageInfo_VoteResponse.Size(m)
}
func (m *VoteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VoteResponse proto.InternalMessageInfo

func (m *VoteResponse) GetRsp() *Response {
	if m != nil {
		return m.Rsp
	}
	return nil
}

func (m *VoteResponse) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

func init() {
	proto.RegisterType((*VoteRequest)(nil), "service.VoteRequest")
	proto.RegisterType((*Response)(nil), "service.Response")
	proto.RegisterType((*VoteResponse)(nil), "service.VoteResponse")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x41, 0x0b, 0x82, 0x30,
	0x14, 0xc7, 0x31, 0xcd, 0xec, 0x59, 0x50, 0xa3, 0x42, 0x3a, 0x89, 0x5d, 0x3c, 0x49, 0xe8, 0xa9,
	0xa3, 0xc7, 0xe8, 0xb6, 0xa0, 0x0f, 0x90, 0x3e, 0x42, 0x48, 0xb7, 0xb6, 0xd9, 0xe7, 0x8f, 0xcd,
	0x19, 0xd1, 0xed, 0xff, 0xdf, 0xe3, 0xf7, 0x7b, 0xdb, 0x60, 0x29, 0x51, 0xbc, 0x9b, 0x0a, 0x33,
	0x2e, 0x98, 0x62, 0x64, 0x66, 0x6b, 0x72, 0x82, 0xf0, 0xc6, 0x14, 0x52, 0x7c, 0xf5, 0x28, 0x15,
	0x21, 0xe0, 0x29, 0x14, 0x6d, 0xe4, 0xc4, 0x4e, 0x3a, 0xa5, 0x26, 0x93, 0x1d, 0xf8, 0x1d, 0xab,
	0xf1, 0x5c, 0x47, 0x93, 0xd8, 0x49, 0xe7, 0xd4, 0xb6, 0xe4, 0x08, 0x01, 0x45, 0xc9, 0x59, 0x27,
	0x51, 0x73, 0x15, 0xab, 0x71, 0xe4, 0x74, 0x26, 0x2b, 0x70, 0x5b, 0xf9, 0xb0, 0x90, 0x8e, 0xc9,
	0x05, 0x16, 0xc3, 0x32, 0x4b, 0x1d, 0xc0, 0x15, 0x92, 0x1b, 0x28, 0xcc, 0xd7, 0xd9, 0x78, 0xc5,
	0x71, 0x4e, 0xf5, 0x54, 0xaf, 0x17, 0x28, 0xfb, 0xa7, 0x32, 0xa6, 0x80, 0xda, 0x96, 0x97, 0x00,
	0x25, 0x6f, 0xae, 0x03, 0x43, 0x0a, 0xf0, 0xb4, 0x9a, 0x6c, 0xbe, 0x96, 0x9f, 0x67, 0xed, 0xb7,
	0x7f, 0xa7, 0x83, 0xff, 0xee, 0x9b, 0xcf, 0x28, 0x3e, 0x01, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x89,
	0xb7, 0x76, 0x1d, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ApiServiceClient is the client API for ApiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApiServiceClient interface {
	// 定义方法
	Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error)
}

type apiServiceClient struct {
	cc *grpc.ClientConn
}

func NewApiServiceClient(cc *grpc.ClientConn) ApiServiceClient {
	return &apiServiceClient{cc}
}

func (c *apiServiceClient) Vote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error) {
	out := new(VoteResponse)
	err := c.cc.Invoke(ctx, "/service.ApiService/Vote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApiServiceServer is the server API for ApiService service.
type ApiServiceServer interface {
	// 定义方法
	Vote(context.Context, *VoteRequest) (*VoteResponse, error)
}

// UnimplementedApiServiceServer can be embedded to have forward compatible implementations.
type UnimplementedApiServiceServer struct {
}

func (*UnimplementedApiServiceServer) Vote(ctx context.Context, req *VoteRequest) (*VoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}

func RegisterApiServiceServer(s *grpc.Server, srv ApiServiceServer) {
	s.RegisterService(&_ApiService_serviceDesc, srv)
}

func _ApiService_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServiceServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ApiService/Vote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServiceServer).Vote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApiService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.ApiService",
	HandlerType: (*ApiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Vote",
			Handler:    _ApiService_Vote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
