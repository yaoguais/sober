// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kv.proto

package kvpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{0}
}
func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (dst *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(dst, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type GetResponse struct {
	Value                string   `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{1}
}
func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (dst *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(dst, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type SetRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Velue                string   `protobuf:"bytes,2,opt,name=velue" json:"velue,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetRequest) Reset()         { *m = SetRequest{} }
func (m *SetRequest) String() string { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()    {}
func (*SetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{2}
}
func (m *SetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetRequest.Unmarshal(m, b)
}
func (m *SetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetRequest.Marshal(b, m, deterministic)
}
func (dst *SetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetRequest.Merge(dst, src)
}
func (m *SetRequest) XXX_Size() int {
	return xxx_messageInfo_SetRequest.Size(m)
}
func (m *SetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetRequest proto.InternalMessageInfo

func (m *SetRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SetRequest) GetVelue() string {
	if m != nil {
		return m.Velue
	}
	return ""
}

type SetResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetResponse) Reset()         { *m = SetResponse{} }
func (m *SetResponse) String() string { return proto.CompactTextString(m) }
func (*SetResponse) ProtoMessage()    {}
func (*SetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{3}
}
func (m *SetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetResponse.Unmarshal(m, b)
}
func (m *SetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetResponse.Marshal(b, m, deterministic)
}
func (dst *SetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetResponse.Merge(dst, src)
}
func (m *SetResponse) XXX_Size() int {
	return xxx_messageInfo_SetResponse.Size(m)
}
func (m *SetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetResponse proto.InternalMessageInfo

type WatchRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WatchRequest) Reset()         { *m = WatchRequest{} }
func (m *WatchRequest) String() string { return proto.CompactTextString(m) }
func (*WatchRequest) ProtoMessage()    {}
func (*WatchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{4}
}
func (m *WatchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WatchRequest.Unmarshal(m, b)
}
func (m *WatchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WatchRequest.Marshal(b, m, deterministic)
}
func (dst *WatchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WatchRequest.Merge(dst, src)
}
func (m *WatchRequest) XXX_Size() int {
	return xxx_messageInfo_WatchRequest.Size(m)
}
func (m *WatchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WatchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WatchRequest proto.InternalMessageInfo

func (m *WatchRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type Event struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{5}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

type FeedbackRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Error                bool     `protobuf:"varint,2,opt,name=error" json:"error,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeedbackRequest) Reset()         { *m = FeedbackRequest{} }
func (m *FeedbackRequest) String() string { return proto.CompactTextString(m) }
func (*FeedbackRequest) ProtoMessage()    {}
func (*FeedbackRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{6}
}
func (m *FeedbackRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeedbackRequest.Unmarshal(m, b)
}
func (m *FeedbackRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeedbackRequest.Marshal(b, m, deterministic)
}
func (dst *FeedbackRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeedbackRequest.Merge(dst, src)
}
func (m *FeedbackRequest) XXX_Size() int {
	return xxx_messageInfo_FeedbackRequest.Size(m)
}
func (m *FeedbackRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FeedbackRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FeedbackRequest proto.InternalMessageInfo

func (m *FeedbackRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *FeedbackRequest) GetError() bool {
	if m != nil {
		return m.Error
	}
	return false
}

func (m *FeedbackRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type FeedbackResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeedbackResponse) Reset()         { *m = FeedbackResponse{} }
func (m *FeedbackResponse) String() string { return proto.CompactTextString(m) }
func (*FeedbackResponse) ProtoMessage()    {}
func (*FeedbackResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_kv_66c22fc9e4b82bbb, []int{7}
}
func (m *FeedbackResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeedbackResponse.Unmarshal(m, b)
}
func (m *FeedbackResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeedbackResponse.Marshal(b, m, deterministic)
}
func (dst *FeedbackResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeedbackResponse.Merge(dst, src)
}
func (m *FeedbackResponse) XXX_Size() int {
	return xxx_messageInfo_FeedbackResponse.Size(m)
}
func (m *FeedbackResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FeedbackResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FeedbackResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetRequest)(nil), "kvpb.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "kvpb.GetResponse")
	proto.RegisterType((*SetRequest)(nil), "kvpb.SetRequest")
	proto.RegisterType((*SetResponse)(nil), "kvpb.SetResponse")
	proto.RegisterType((*WatchRequest)(nil), "kvpb.WatchRequest")
	proto.RegisterType((*Event)(nil), "kvpb.Event")
	proto.RegisterType((*FeedbackRequest)(nil), "kvpb.FeedbackRequest")
	proto.RegisterType((*FeedbackResponse)(nil), "kvpb.FeedbackResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// KVClient is the client API for KV service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KVClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (KV_WatchClient, error)
	Feedback(ctx context.Context, in *FeedbackRequest, opts ...grpc.CallOption) (*FeedbackResponse, error)
}

type kVClient struct {
	cc *grpc.ClientConn
}

func NewKVClient(cc *grpc.ClientConn) KVClient {
	return &kVClient{cc}
}

func (c *kVClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/kvpb.KV/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/kvpb.KV/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (KV_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_KV_serviceDesc.Streams[0], "/kvpb.KV/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &kVWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KV_WatchClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type kVWatchClient struct {
	grpc.ClientStream
}

func (x *kVWatchClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *kVClient) Feedback(ctx context.Context, in *FeedbackRequest, opts ...grpc.CallOption) (*FeedbackResponse, error) {
	out := new(FeedbackResponse)
	err := c.cc.Invoke(ctx, "/kvpb.KV/Feedback", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVServer is the server API for KV service.
type KVServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Watch(*WatchRequest, KV_WatchServer) error
	Feedback(context.Context, *FeedbackRequest) (*FeedbackResponse, error)
}

func RegisterKVServer(s *grpc.Server, srv KVServer) {
	s.RegisterService(&_KV_serviceDesc, srv)
}

func _KV_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvpb.KV/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KV_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvpb.KV/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KV_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KVServer).Watch(m, &kVWatchServer{stream})
}

type KV_WatchServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type kVWatchServer struct {
	grpc.ServerStream
}

func (x *kVWatchServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func _KV_Feedback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServer).Feedback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvpb.KV/Feedback",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).Feedback(ctx, req.(*FeedbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _KV_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kvpb.KV",
	HandlerType: (*KVServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _KV_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _KV_Set_Handler,
		},
		{
			MethodName: "Feedback",
			Handler:    _KV_Feedback_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _KV_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "kv.proto",
}

func init() { proto.RegisterFile("kv.proto", fileDescriptor_kv_66c22fc9e4b82bbb) }

var fileDescriptor_kv_66c22fc9e4b82bbb = []byte{
	// 267 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x4f, 0x4f, 0x83, 0x40,
	0x10, 0xc5, 0x4b, 0x11, 0x8b, 0x0f, 0x8d, 0x38, 0x51, 0x43, 0x38, 0x98, 0x66, 0xbd, 0x78, 0x30,
	0xc4, 0xa8, 0x37, 0xcf, 0xda, 0x83, 0x37, 0x36, 0xd1, 0x33, 0xd4, 0x89, 0x26, 0x68, 0x41, 0xd8,
	0x92, 0xf8, 0x19, 0xfd, 0x52, 0x66, 0x77, 0x21, 0x34, 0x98, 0xf6, 0xb6, 0x6f, 0xe6, 0x37, 0xff,
	0xde, 0xc2, 0x2f, 0xda, 0xa4, 0xaa, 0x4b, 0x55, 0xd2, 0x5e, 0xd1, 0x56, 0xb9, 0xb8, 0x00, 0x16,
	0xac, 0x52, 0xfe, 0x5e, 0x73, 0xa3, 0x28, 0x84, 0x5b, 0xf0, 0x4f, 0xe4, 0xcc, 0x9d, 0xab, 0x83,
	0x54, 0x3f, 0xc5, 0x25, 0x02, 0x93, 0x6f, 0xaa, 0x72, 0xd5, 0x30, 0x9d, 0xc2, 0x6b, 0xb3, 0xcf,
	0x35, 0x77, 0x88, 0x15, 0xe2, 0x1e, 0x90, 0x3b, 0x9a, 0x98, 0x2a, 0xd6, 0x55, 0xd3, 0xae, 0x4a,
	0x0b, 0x71, 0x84, 0x40, 0x0e, 0xad, 0xc5, 0x1c, 0x87, 0xaf, 0x99, 0x5a, 0x7e, 0x6c, 0xdf, 0x65,
	0x06, 0xef, 0xb1, 0xe5, 0x95, 0x12, 0x12, 0xc7, 0x4f, 0xcc, 0x6f, 0x79, 0xb6, 0x2c, 0x76, 0x0e,
	0xe5, 0xba, 0x2e, 0x6b, 0x33, 0xd4, 0x4f, 0xad, 0xa0, 0x08, 0xb3, 0x2f, 0x6e, 0x9a, 0xec, 0x9d,
	0x23, 0xd7, 0xb0, 0xbd, 0x14, 0x84, 0x70, 0x68, 0x6a, 0x77, 0xba, 0xfd, 0x75, 0x30, 0x7d, 0x7e,
	0xa1, 0x6b, 0xb8, 0x0b, 0x56, 0x14, 0x26, 0xda, 0xb2, 0x64, 0xf0, 0x2b, 0x3e, 0xd9, 0x88, 0x74,
	0x67, 0x4c, 0x34, 0x2d, 0x07, 0x5a, 0xfe, 0xa3, 0xe5, 0x88, 0xf6, 0xcc, 0xd9, 0x44, 0x36, 0xbb,
	0xe9, 0x41, 0x1c, 0xd8, 0x98, 0xbd, 0x7a, 0x72, 0xe3, 0xd0, 0x03, 0xfc, 0x7e, 0x49, 0x3a, 0xb3,
	0xc9, 0x91, 0x13, 0xf1, 0xf9, 0x38, 0xdc, 0x8f, 0xca, 0xf7, 0xcd, 0xc7, 0xdf, 0xfd, 0x05, 0x00,
	0x00, 0xff, 0xff, 0x1c, 0x8e, 0x89, 0x09, 0x04, 0x02, 0x00, 0x00,
}
