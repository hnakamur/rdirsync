// Code generated by protoc-gen-go.
// source: rdirsync.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	rdirsync.proto

It has these top-level messages:
	StatRequest
	FetchFileRequest
	FileChunk
	ReadDirRequest
	FileInfos
	FileInfo
	UnixTime
	ChmodRequest
	ChtimesRequest
	EnsureNotExistRequest
	Empty
	EnsureDirExistsRequest
	SendFileRequest
*/
package pb

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
	Name    string    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Size    int64     `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Mode    int32     `protobuf:"varint,3,opt,name=mode" json:"mode,omitempty"`
	ModTime *UnixTime `protobuf:"bytes,4,opt,name=modTime" json:"modTime,omitempty"`
	Owner   string    `protobuf:"bytes,5,opt,name=owner" json:"owner,omitempty"`
	Group   string    `protobuf:"bytes,6,opt,name=group" json:"group,omitempty"`
}

func (m *FileInfo) Reset()                    { *m = FileInfo{} }
func (m *FileInfo) String() string            { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()               {}
func (*FileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *FileInfo) GetModTime() *UnixTime {
	if m != nil {
		return m.ModTime
	}
	return nil
}

type UnixTime struct {
	Second     int64 `protobuf:"varint,1,opt,name=second" json:"second,omitempty"`
	NanoSecond int64 `protobuf:"varint,2,opt,name=nanoSecond" json:"nanoSecond,omitempty"`
}

func (m *UnixTime) Reset()                    { *m = UnixTime{} }
func (m *UnixTime) String() string            { return proto.CompactTextString(m) }
func (*UnixTime) ProtoMessage()               {}
func (*UnixTime) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type ChmodRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Mode int32  `protobuf:"varint,2,opt,name=mode" json:"mode,omitempty"`
}

func (m *ChmodRequest) Reset()                    { *m = ChmodRequest{} }
func (m *ChmodRequest) String() string            { return proto.CompactTextString(m) }
func (*ChmodRequest) ProtoMessage()               {}
func (*ChmodRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type ChtimesRequest struct {
	Path  string    `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Atime *UnixTime `protobuf:"bytes,2,opt,name=atime" json:"atime,omitempty"`
	Mtime *UnixTime `protobuf:"bytes,3,opt,name=mtime" json:"mtime,omitempty"`
}

func (m *ChtimesRequest) Reset()                    { *m = ChtimesRequest{} }
func (m *ChtimesRequest) String() string            { return proto.CompactTextString(m) }
func (*ChtimesRequest) ProtoMessage()               {}
func (*ChtimesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ChtimesRequest) GetAtime() *UnixTime {
	if m != nil {
		return m.Atime
	}
	return nil
}

func (m *ChtimesRequest) GetMtime() *UnixTime {
	if m != nil {
		return m.Mtime
	}
	return nil
}

type EnsureNotExistRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *EnsureNotExistRequest) Reset()                    { *m = EnsureNotExistRequest{} }
func (m *EnsureNotExistRequest) String() string            { return proto.CompactTextString(m) }
func (*EnsureNotExistRequest) ProtoMessage()               {}
func (*EnsureNotExistRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type EnsureDirExistsRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *EnsureDirExistsRequest) Reset()                    { *m = EnsureDirExistsRequest{} }
func (m *EnsureDirExistsRequest) String() string            { return proto.CompactTextString(m) }
func (*EnsureDirExistsRequest) ProtoMessage()               {}
func (*EnsureDirExistsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type SendFileRequest struct {
	Path  string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Chunk []byte `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (m *SendFileRequest) Reset()                    { *m = SendFileRequest{} }
func (m *SendFileRequest) String() string            { return proto.CompactTextString(m) }
func (*SendFileRequest) ProtoMessage()               {}
func (*SendFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func init() {
	proto.RegisterType((*StatRequest)(nil), "pb.StatRequest")
	proto.RegisterType((*FetchFileRequest)(nil), "pb.FetchFileRequest")
	proto.RegisterType((*FileChunk)(nil), "pb.FileChunk")
	proto.RegisterType((*ReadDirRequest)(nil), "pb.ReadDirRequest")
	proto.RegisterType((*FileInfos)(nil), "pb.FileInfos")
	proto.RegisterType((*FileInfo)(nil), "pb.FileInfo")
	proto.RegisterType((*UnixTime)(nil), "pb.UnixTime")
	proto.RegisterType((*ChmodRequest)(nil), "pb.ChmodRequest")
	proto.RegisterType((*ChtimesRequest)(nil), "pb.ChtimesRequest")
	proto.RegisterType((*EnsureNotExistRequest)(nil), "pb.EnsureNotExistRequest")
	proto.RegisterType((*Empty)(nil), "pb.Empty")
	proto.RegisterType((*EnsureDirExistsRequest)(nil), "pb.EnsureDirExistsRequest")
	proto.RegisterType((*SendFileRequest)(nil), "pb.SendFileRequest")
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
	Chmod(ctx context.Context, in *ChmodRequest, opts ...grpc.CallOption) (*Empty, error)
	Chtimes(ctx context.Context, in *ChtimesRequest, opts ...grpc.CallOption) (*Empty, error)
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
	err := grpc.Invoke(ctx, "/pb.RDirSync/Stat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDirSyncClient) ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (RDirSync_ReadDirClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[0], c.cc, "/pb.RDirSync/ReadDir", opts...)
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
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[1], c.cc, "/pb.RDirSync/FetchFile", opts...)
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
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[2], c.cc, "/pb.RDirSync/SendFile", opts...)
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

func (c *rDirSyncClient) Chmod(ctx context.Context, in *ChmodRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/pb.RDirSync/Chmod", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDirSyncClient) Chtimes(ctx context.Context, in *ChtimesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/pb.RDirSync/Chtimes", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDirSyncClient) EnsureDirExists(ctx context.Context, in *EnsureDirExistsRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/pb.RDirSync/EnsureDirExists", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDirSyncClient) EnsureNotExist(ctx context.Context, in *EnsureNotExistRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/pb.RDirSync/EnsureNotExist", in, out, c.cc, opts...)
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
	Chmod(context.Context, *ChmodRequest) (*Empty, error)
	Chtimes(context.Context, *ChtimesRequest) (*Empty, error)
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
		FullMethod: "/pb.RDirSync/Stat",
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

func _RDirSync_Chmod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChmodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDirSyncServer).Chmod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RDirSync/Chmod",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).Chmod(ctx, req.(*ChmodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDirSync_Chtimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChtimesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDirSyncServer).Chtimes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RDirSync/Chtimes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).Chtimes(ctx, req.(*ChtimesRequest))
	}
	return interceptor(ctx, in, info, handler)
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
		FullMethod: "/pb.RDirSync/EnsureDirExists",
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
		FullMethod: "/pb.RDirSync/EnsureNotExist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).EnsureNotExist(ctx, req.(*EnsureNotExistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RDirSync_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.RDirSync",
	HandlerType: (*RDirSyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Stat",
			Handler:    _RDirSync_Stat_Handler,
		},
		{
			MethodName: "Chmod",
			Handler:    _RDirSync_Chmod_Handler,
		},
		{
			MethodName: "Chtimes",
			Handler:    _RDirSync_Chtimes_Handler,
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
	// 574 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x54, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0xad, 0xeb, 0x3a, 0x97, 0x49, 0x9a, 0x54, 0x43, 0xa9, 0x4c, 0x1e, 0x50, 0xba, 0x48, 0x25,
	0x02, 0x64, 0xaa, 0x82, 0x2a, 0x04, 0x2f, 0x28, 0x37, 0x89, 0x07, 0x50, 0xe5, 0xc0, 0x07, 0xf8,
	0xb2, 0xad, 0x57, 0xad, 0x77, 0x8d, 0xbd, 0x16, 0x0d, 0xff, 0xc1, 0x0b, 0x5f, 0x8b, 0x76, 0x1d,
	0x37, 0x76, 0x54, 0x99, 0xb7, 0xd9, 0x73, 0xe6, 0x78, 0x66, 0xbc, 0x67, 0x16, 0x06, 0x69, 0xc8,
	0xd2, 0x6c, 0xcd, 0x03, 0x27, 0x49, 0x85, 0x14, 0xb8, 0x9f, 0xf8, 0xe4, 0x14, 0x7a, 0x2b, 0xe9,
	0x49, 0x97, 0xfe, 0xcc, 0x69, 0x26, 0x11, 0xe1, 0x20, 0xf1, 0x64, 0x64, 0x1b, 0x63, 0x63, 0xd2,
	0x75, 0x75, 0x4c, 0x3e, 0xc3, 0xd1, 0x92, 0xca, 0x20, 0x5a, 0xb2, 0x3b, 0xda, 0x90, 0x87, 0x36,
	0xb4, 0xfd, 0xfc, 0x7a, 0xc5, 0x7e, 0x53, 0x7b, 0x7f, 0x6c, 0x4c, 0x2c, 0xb7, 0x3c, 0x92, 0x53,
	0xe8, 0x2a, 0xf1, 0x2c, 0xca, 0xf9, 0x2d, 0x1e, 0x83, 0x15, 0xa8, 0x40, 0x6b, 0xfb, 0x6e, 0x71,
	0x20, 0x4b, 0x18, 0xb8, 0xd4, 0x0b, 0xe7, 0x2c, 0x6d, 0x2a, 0x31, 0x86, 0x9e, 0x27, 0xbf, 0x8a,
	0x4c, 0xce, 0x44, 0xce, 0xe5, 0xa6, 0x4c, 0x15, 0x22, 0x6f, 0x8b, 0x52, 0x5f, 0xf8, 0xb5, 0xc8,
	0x90, 0x80, 0xc5, 0x54, 0x60, 0x1b, 0x63, 0x73, 0xd2, 0xbb, 0xe8, 0x3b, 0x89, 0xef, 0x94, 0xac,
	0x5b, 0x50, 0xe4, 0xaf, 0x01, 0x9d, 0x12, 0x53, 0x35, 0xb9, 0x17, 0xd3, 0xb2, 0xa6, 0x8a, 0x15,
	0x96, 0x95, 0x33, 0x99, 0xae, 0x8e, 0x15, 0x16, 0x8b, 0x90, 0xda, 0xa6, 0x6e, 0x40, 0xc7, 0x78,
	0x06, 0xed, 0x58, 0x84, 0xdf, 0x59, 0x4c, 0xed, 0x83, 0xb1, 0x51, 0x96, 0xfb, 0xc1, 0xd9, 0xbd,
	0xc2, 0xdc, 0x92, 0x54, 0xf3, 0x8b, 0x5f, 0x9c, 0xa6, 0xb6, 0xa5, 0x8b, 0x14, 0x07, 0x85, 0xde,
	0xa4, 0x22, 0x4f, 0xec, 0x56, 0x81, 0xea, 0x03, 0x99, 0x42, 0xa7, 0xfc, 0x00, 0x9e, 0x40, 0x2b,
	0xa3, 0x81, 0xe0, 0xa1, 0xee, 0xce, 0x74, 0x37, 0x27, 0x7c, 0x0e, 0xc0, 0x3d, 0x2e, 0x56, 0x05,
	0x57, 0x74, 0x59, 0x41, 0xc8, 0x25, 0xf4, 0x67, 0x51, 0x2c, 0xc2, 0xa6, 0xff, 0x5a, 0xce, 0xb3,
	0xbf, 0x9d, 0x87, 0xdc, 0xc1, 0x60, 0x16, 0x49, 0x16, 0xd3, 0xac, 0x49, 0x49, 0xc0, 0xf2, 0x54,
	0x92, 0x96, 0xee, 0xce, 0x5c, 0x50, 0x2a, 0x27, 0xd6, 0x39, 0xe6, 0x63, 0x39, 0x9a, 0x22, 0xaf,
	0xe1, 0xe9, 0x82, 0x67, 0x79, 0x4a, 0xbf, 0x09, 0xb9, 0xb8, 0x67, 0x59, 0xa3, 0x23, 0xdb, 0x60,
	0x2d, 0xe2, 0x44, 0xae, 0xc9, 0x1b, 0x38, 0x29, 0x54, 0x73, 0x96, 0x6a, 0x55, 0x53, 0xaf, 0xe4,
	0x13, 0x0c, 0x57, 0x94, 0x87, 0xff, 0xf3, 0xf1, 0x83, 0x41, 0xcd, 0x8a, 0x41, 0x2f, 0xfe, 0x98,
	0xd0, 0x71, 0xe7, 0x2c, 0x5d, 0xad, 0x79, 0x80, 0x2f, 0xe1, 0x40, 0x6d, 0x0d, 0x0e, 0xd5, 0x28,
	0x95, 0xfd, 0x19, 0xd5, 0x2c, 0x46, 0xf6, 0xf0, 0x1c, 0xda, 0x1b, 0x5b, 0x23, 0x2a, 0xaa, 0xee,
	0xf1, 0xd1, 0x61, 0x35, 0x3d, 0x23, 0x7b, 0xe7, 0x06, 0xbe, 0x87, 0xee, 0xc3, 0xb6, 0xe1, 0xb1,
	0xe6, 0x77, 0x96, 0x6f, 0xab, 0xd2, 0x0b, 0xa5, 0x55, 0x0e, 0x74, 0xca, 0xd1, 0xf0, 0x89, 0x6e,
	0xaa, 0x3e, 0xe8, 0xa8, 0xab, 0xc0, 0xe2, 0xa7, 0xed, 0x4d, 0x0c, 0x3c, 0x03, 0x4b, 0x9b, 0x02,
	0x8f, 0x14, 0x5e, 0xf5, 0x47, 0x2d, 0x13, 0x5f, 0x41, 0x7b, 0x63, 0x82, 0xa2, 0xff, 0xba, 0x23,
	0xea, 0xb9, 0x1f, 0x61, 0xb8, 0x73, 0x19, 0x38, 0xd2, 0xfc, 0xa3, 0x37, 0x54, 0xd7, 0x7e, 0x80,
	0x41, 0xfd, 0xfa, 0xf1, 0xd9, 0x56, 0xba, 0x63, 0x89, 0x9a, 0x72, 0x7a, 0x09, 0x2f, 0x02, 0x11,
	0x3b, 0x37, 0x4c, 0x46, 0xb9, 0xef, 0xa8, 0x30, 0xe2, 0xde, 0xad, 0x17, 0xe7, 0xa9, 0xb3, 0x7d,
	0xee, 0xfc, 0xe9, 0x61, 0x79, 0x77, 0x57, 0xea, 0xe9, 0xbb, 0x32, 0xfc, 0x96, 0x7e, 0x03, 0xdf,
	0xfd, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x15, 0x81, 0xf3, 0xf7, 0x15, 0x05, 0x00, 0x00,
}
