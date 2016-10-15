// Code generated by protoc-gen-go.
// source: rdirsync.proto
// DO NOT EDIT!

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	rdirsync.proto

It has these top-level messages:
	StatRequest
	FetchFileRequest
	FileChunk
	ReadDirRequest
	FileInfos
	FileInfo
	EnsureNotExistRequest
	Empty
	EnsureDirExistsRequest
	SendFileRequest
*/
package rpc

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

type StatRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *StatRequest) Reset()                    { *m = StatRequest{} }
func (m *StatRequest) String() string            { return proto.CompactTextString(m) }
func (*StatRequest) ProtoMessage()               {}
func (*StatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type FetchFileRequest struct {
	Path    string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	BufSize int32  `protobuf:"varint,2,opt,name=bufSize" json:"bufSize,omitempty"`
}

func (m *FetchFileRequest) Reset()                    { *m = FetchFileRequest{} }
func (m *FetchFileRequest) String() string            { return proto.CompactTextString(m) }
func (*FetchFileRequest) ProtoMessage()               {}
func (*FetchFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type FileChunk struct {
	Chunk []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (m *FileChunk) Reset()                    { *m = FileChunk{} }
func (m *FileChunk) String() string            { return proto.CompactTextString(m) }
func (*FileChunk) ProtoMessage()               {}
func (*FileChunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type ReadDirRequest struct {
	Path        string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	AtMostCount int32  `protobuf:"varint,2,opt,name=atMostCount" json:"atMostCount,omitempty"`
}

func (m *ReadDirRequest) Reset()                    { *m = ReadDirRequest{} }
func (m *ReadDirRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadDirRequest) ProtoMessage()               {}
func (*ReadDirRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type FileInfos struct {
	Infos []*FileInfo `protobuf:"bytes,1,rep,name=infos" json:"infos,omitempty"`
}

func (m *FileInfos) Reset()                    { *m = FileInfos{} }
func (m *FileInfos) String() string            { return proto.CompactTextString(m) }
func (*FileInfos) ProtoMessage()               {}
func (*FileInfos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *FileInfos) GetInfos() []*FileInfo {
	if m != nil {
		return m.Infos
	}
	return nil
}

type FileInfo struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Size    int64  `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Mode    int32  `protobuf:"varint,3,opt,name=mode" json:"mode,omitempty"`
	ModTime int64  `protobuf:"varint,4,opt,name=modTime" json:"modTime,omitempty"`
}

func (m *FileInfo) Reset()                    { *m = FileInfo{} }
func (m *FileInfo) String() string            { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()               {}
func (*FileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type EnsureNotExistRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *EnsureNotExistRequest) Reset()                    { *m = EnsureNotExistRequest{} }
func (m *EnsureNotExistRequest) String() string            { return proto.CompactTextString(m) }
func (*EnsureNotExistRequest) ProtoMessage()               {}
func (*EnsureNotExistRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type EnsureDirExistsRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *EnsureDirExistsRequest) Reset()                    { *m = EnsureDirExistsRequest{} }
func (m *EnsureDirExistsRequest) String() string            { return proto.CompactTextString(m) }
func (*EnsureDirExistsRequest) ProtoMessage()               {}
func (*EnsureDirExistsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type SendFileRequest struct {
	Path  string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Mode  int32  `protobuf:"varint,2,opt,name=mode" json:"mode,omitempty"`
	Chunk []byte `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (m *SendFileRequest) Reset()                    { *m = SendFileRequest{} }
func (m *SendFileRequest) String() string            { return proto.CompactTextString(m) }
func (*SendFileRequest) ProtoMessage()               {}
func (*SendFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func init() {
	proto.RegisterType((*StatRequest)(nil), "rpc.StatRequest")
	proto.RegisterType((*FetchFileRequest)(nil), "rpc.FetchFileRequest")
	proto.RegisterType((*FileChunk)(nil), "rpc.FileChunk")
	proto.RegisterType((*ReadDirRequest)(nil), "rpc.ReadDirRequest")
	proto.RegisterType((*FileInfos)(nil), "rpc.FileInfos")
	proto.RegisterType((*FileInfo)(nil), "rpc.FileInfo")
	proto.RegisterType((*EnsureNotExistRequest)(nil), "rpc.EnsureNotExistRequest")
	proto.RegisterType((*Empty)(nil), "rpc.Empty")
	proto.RegisterType((*EnsureDirExistsRequest)(nil), "rpc.EnsureDirExistsRequest")
	proto.RegisterType((*SendFileRequest)(nil), "rpc.SendFileRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for RDirSync service

type RDirSyncClient interface {
	Stat(ctx context.Context, in *StatRequest, opts ...grpc.CallOption) (*FileInfo, error)
	ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (RDirSync_ReadDirClient, error)
	FetchFile(ctx context.Context, in *FetchFileRequest, opts ...grpc.CallOption) (RDirSync_FetchFileClient, error)
	SendFile(ctx context.Context, opts ...grpc.CallOption) (RDirSync_SendFileClient, error)
	EnsureDirExists(ctx context.Context, in *EnsureDirExistsRequest, opts ...grpc.CallOption) (*Empty, error)
	EnsureNotExist(ctx context.Context, in *EnsureNotExistRequest, opts ...grpc.CallOption) (*Empty, error)
}

type rDirSyncClient struct {
	cc *grpc.ClientConn
}

func NewRDirSyncClient(cc *grpc.ClientConn) RDirSyncClient {
	return &rDirSyncClient{cc}
}

func (c *rDirSyncClient) Stat(ctx context.Context, in *StatRequest, opts ...grpc.CallOption) (*FileInfo, error) {
	out := new(FileInfo)
	err := grpc.Invoke(ctx, "/rpc.RDirSync/Stat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDirSyncClient) ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (RDirSync_ReadDirClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[0], c.cc, "/rpc.RDirSync/ReadDir", opts...)
	if err != nil {
		return nil, err
	}
	x := &rDirSyncReadDirClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RDirSync_ReadDirClient interface {
	Recv() (*FileInfos, error)
	grpc.ClientStream
}

type rDirSyncReadDirClient struct {
	grpc.ClientStream
}

func (x *rDirSyncReadDirClient) Recv() (*FileInfos, error) {
	m := new(FileInfos)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rDirSyncClient) FetchFile(ctx context.Context, in *FetchFileRequest, opts ...grpc.CallOption) (RDirSync_FetchFileClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[1], c.cc, "/rpc.RDirSync/FetchFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &rDirSyncFetchFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RDirSync_FetchFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type rDirSyncFetchFileClient struct {
	grpc.ClientStream
}

func (x *rDirSyncFetchFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rDirSyncClient) SendFile(ctx context.Context, opts ...grpc.CallOption) (RDirSync_SendFileClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[2], c.cc, "/rpc.RDirSync/SendFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &rDirSyncSendFileClient{stream}
	return x, nil
}

type RDirSync_SendFileClient interface {
	Send(*SendFileRequest) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type rDirSyncSendFileClient struct {
	grpc.ClientStream
}

func (x *rDirSyncSendFileClient) Send(m *SendFileRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *rDirSyncSendFileClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rDirSyncClient) EnsureDirExists(ctx context.Context, in *EnsureDirExistsRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/rpc.RDirSync/EnsureDirExists", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDirSyncClient) EnsureNotExist(ctx context.Context, in *EnsureNotExistRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/rpc.RDirSync/EnsureNotExist", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RDirSync service

type RDirSyncServer interface {
	Stat(context.Context, *StatRequest) (*FileInfo, error)
	ReadDir(*ReadDirRequest, RDirSync_ReadDirServer) error
	FetchFile(*FetchFileRequest, RDirSync_FetchFileServer) error
	SendFile(RDirSync_SendFileServer) error
	EnsureDirExists(context.Context, *EnsureDirExistsRequest) (*Empty, error)
	EnsureNotExist(context.Context, *EnsureNotExistRequest) (*Empty, error)
}

func RegisterRDirSyncServer(s *grpc.Server, srv RDirSyncServer) {
	s.RegisterService(&_RDirSync_serviceDesc, srv)
}

func _RDirSync_Stat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDirSyncServer).Stat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.RDirSync/Stat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).Stat(ctx, req.(*StatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDirSync_ReadDir_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReadDirRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RDirSyncServer).ReadDir(m, &rDirSyncReadDirServer{stream})
}

type RDirSync_ReadDirServer interface {
	Send(*FileInfos) error
	grpc.ServerStream
}

type rDirSyncReadDirServer struct {
	grpc.ServerStream
}

func (x *rDirSyncReadDirServer) Send(m *FileInfos) error {
	return x.ServerStream.SendMsg(m)
}

func _RDirSync_FetchFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FetchFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RDirSyncServer).FetchFile(m, &rDirSyncFetchFileServer{stream})
}

type RDirSync_FetchFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type rDirSyncFetchFileServer struct {
	grpc.ServerStream
}

func (x *rDirSyncFetchFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

func _RDirSync_SendFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RDirSyncServer).SendFile(&rDirSyncSendFileServer{stream})
}

type RDirSync_SendFileServer interface {
	SendAndClose(*Empty) error
	Recv() (*SendFileRequest, error)
	grpc.ServerStream
}

type rDirSyncSendFileServer struct {
	grpc.ServerStream
}

func (x *rDirSyncSendFileServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *rDirSyncSendFileServer) Recv() (*SendFileRequest, error) {
	m := new(SendFileRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RDirSync_EnsureDirExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnsureDirExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDirSyncServer).EnsureDirExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.RDirSync/EnsureDirExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).EnsureDirExists(ctx, req.(*EnsureDirExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDirSync_EnsureNotExist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnsureNotExistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDirSyncServer).EnsureNotExist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.RDirSync/EnsureNotExist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).EnsureNotExist(ctx, req.(*EnsureNotExistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RDirSync_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.RDirSync",
	HandlerType: (*RDirSyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Stat",
			Handler:    _RDirSync_Stat_Handler,
		},
		{
			MethodName: "EnsureDirExists",
			Handler:    _RDirSync_EnsureDirExists_Handler,
		},
		{
			MethodName: "EnsureNotExist",
			Handler:    _RDirSync_EnsureNotExist_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReadDir",
			Handler:       _RDirSync_ReadDir_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "FetchFile",
			Handler:       _RDirSync_FetchFile_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SendFile",
			Handler:       _RDirSync_SendFile_Handler,
			ClientStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("rdirsync.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 455 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x53, 0x51, 0x6f, 0xd3, 0x30,
	0x10, 0x6e, 0x9a, 0x96, 0xb6, 0x57, 0xd6, 0x4d, 0xc7, 0x86, 0xa2, 0xf2, 0xd2, 0x79, 0x2f, 0x45,
	0xa0, 0xa8, 0x2a, 0xd2, 0x1e, 0x10, 0x0f, 0x68, 0x6b, 0x2b, 0xf1, 0x00, 0x4c, 0x29, 0x3f, 0x80,
	0x34, 0xf1, 0xa8, 0x55, 0x62, 0x07, 0xdb, 0x91, 0x28, 0xff, 0x91, 0xff, 0x84, 0xec, 0xd4, 0x55,
	0x52, 0xa1, 0xec, 0xed, 0xf3, 0xf9, 0x3e, 0xdd, 0xdd, 0xf7, 0xdd, 0xc1, 0x48, 0xa6, 0x4c, 0xaa,
	0x3d, 0x4f, 0xc2, 0x5c, 0x0a, 0x2d, 0xd0, 0x97, 0x79, 0x42, 0xae, 0x61, 0xb8, 0xd6, 0xb1, 0x8e,
	0xe8, 0xaf, 0x82, 0x2a, 0x8d, 0x08, 0x9d, 0x3c, 0xd6, 0xdb, 0xc0, 0x9b, 0x78, 0xd3, 0x41, 0x64,
	0x31, 0xf9, 0x08, 0x17, 0x2b, 0xaa, 0x93, 0xed, 0x8a, 0xfd, 0xa4, 0x0d, 0x79, 0x18, 0x40, 0x6f,
	0x53, 0x3c, 0xae, 0xd9, 0x1f, 0x1a, 0xb4, 0x27, 0xde, 0xb4, 0x1b, 0xb9, 0x27, 0xb9, 0x86, 0x81,
	0x21, 0xdf, 0x6f, 0x0b, 0xbe, 0xc3, 0x4b, 0xe8, 0x26, 0x06, 0x58, 0xee, 0xf3, 0xa8, 0x7c, 0x90,
	0x15, 0x8c, 0x22, 0x1a, 0xa7, 0x0b, 0x26, 0x9b, 0x4a, 0x4c, 0x60, 0x18, 0xeb, 0xcf, 0x42, 0xe9,
	0x7b, 0x51, 0x70, 0x7d, 0x28, 0x53, 0x0d, 0x91, 0x59, 0x59, 0xea, 0x13, 0x7f, 0x14, 0x0a, 0x6f,
	0xa0, 0xcb, 0x0c, 0x08, 0xbc, 0x89, 0x3f, 0x1d, 0xce, 0xcf, 0x42, 0x99, 0x27, 0xa1, 0xfb, 0x8e,
	0xca, 0x3f, 0xf2, 0x1d, 0xfa, 0x2e, 0x64, 0x6a, 0xf2, 0x38, 0xa3, 0xae, 0xa6, 0xc1, 0x26, 0xa6,
	0xdc, 0x4c, 0x7e, 0x64, 0xb1, 0x89, 0x65, 0x22, 0xa5, 0x81, 0x6f, 0x1b, 0xb0, 0xd8, 0x8c, 0x9f,
	0x89, 0xf4, 0x1b, 0xcb, 0x68, 0xd0, 0xb1, 0xa9, 0xee, 0x49, 0xde, 0xc0, 0xd5, 0x92, 0xab, 0x42,
	0xd2, 0x2f, 0x42, 0x2f, 0x7f, 0x33, 0xd5, 0xa8, 0x76, 0x0f, 0xba, 0xcb, 0x2c, 0xd7, 0x7b, 0xf2,
	0x16, 0x5e, 0x96, 0xac, 0x05, 0x93, 0x96, 0xa5, 0x9a, 0x68, 0x5f, 0xe1, 0x7c, 0x4d, 0x79, 0xfa,
	0x94, 0x47, 0xae, 0xf1, 0x76, 0xa5, 0xf1, 0xa3, 0x21, 0x7e, 0xc5, 0x90, 0xf9, 0xdf, 0x36, 0xf4,
	0xa3, 0x05, 0x93, 0xeb, 0x3d, 0x4f, 0xf0, 0x35, 0x74, 0xcc, 0x96, 0xe0, 0x85, 0x55, 0xb0, 0xb2,
	0x30, 0xe3, 0xba, 0xa6, 0xa4, 0x85, 0x73, 0xe8, 0x1d, 0x8c, 0xc4, 0x17, 0xf6, 0xaf, 0x6e, 0xeb,
	0x78, 0x54, 0x23, 0x28, 0xd2, 0x9a, 0x79, 0x78, 0x0b, 0x83, 0xe3, 0x86, 0xe1, 0x55, 0x99, 0x70,
	0xb2, 0x71, 0x15, 0x9e, 0x5d, 0x23, 0xcb, 0x9b, 0x41, 0xdf, 0x0d, 0x8d, 0x97, 0x65, 0x6b, 0x75,
	0x0d, 0xc6, 0x60, 0xa3, 0xa5, 0xa0, 0xad, 0xa9, 0x87, 0x1f, 0xe0, 0xfc, 0x44, 0x54, 0x7c, 0x55,
	0xa6, 0xfc, 0x57, 0xea, 0x3a, 0x1f, 0xdf, 0xc3, 0xa8, 0x6e, 0x24, 0x8e, 0x2b, 0xe4, 0x13, 0x77,
	0xeb, 0xdc, 0xbb, 0x5b, 0xb8, 0x11, 0xf2, 0x47, 0xb8, 0x61, 0x7a, 0x53, 0x24, 0x3b, 0xaa, 0xc3,
	0x2d, 0x8f, 0x77, 0x71, 0x56, 0xc8, 0xf0, 0x78, 0x96, 0x32, 0x4f, 0xee, 0xce, 0x9c, 0xe6, 0x0f,
	0xe6, 0x46, 0x1f, 0xbc, 0xcd, 0x33, 0x7b, 0xac, 0xef, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x6e,
	0x19, 0xd2, 0xd9, 0xbe, 0x03, 0x00, 0x00,
}
