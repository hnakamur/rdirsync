// Code generated by protoc-gen-go.
// source: rdirsync.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	rdirsync.proto

It has these top-level messages:
	StatRequest
	ReadDirRequest
	FileInfos
	FileInfo
	FetchFileRequest
	FileChunk
	SendFileRequest
	Empty
	ChownRequest
	ChmodRequest
	ChtimesRequest
	ChangeAttributesRequest
	EnsureNotExistRequest
	EnsureDirExistsRequest
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
	Path               string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	WantsOwnerAndGroup bool   `protobuf:"varint,2,opt,name=wantsOwnerAndGroup" json:"wantsOwnerAndGroup,omitempty"`
}

func (m *StatRequest) Reset()                    { *m = StatRequest{} }
func (m *StatRequest) String() string            { return proto.CompactTextString(m) }
func (*StatRequest) ProtoMessage()               {}
func (*StatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ReadDirRequest struct {
	Path               string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	AtMostCount        int32  `protobuf:"varint,2,opt,name=atMostCount" json:"atMostCount,omitempty"`
	WantsOwnerAndGroup bool   `protobuf:"varint,3,opt,name=wantsOwnerAndGroup" json:"wantsOwnerAndGroup,omitempty"`
}

func (m *ReadDirRequest) Reset()                    { *m = ReadDirRequest{} }
func (m *ReadDirRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadDirRequest) ProtoMessage()               {}
func (*ReadDirRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type FileInfos struct {
	Infos []*FileInfo `protobuf:"bytes,1,rep,name=infos" json:"infos,omitempty"`
}

func (m *FileInfos) Reset()                    { *m = FileInfos{} }
func (m *FileInfos) String() string            { return proto.CompactTextString(m) }
func (*FileInfos) ProtoMessage()               {}
func (*FileInfos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *FileInfos) GetInfos() []*FileInfo {
	if m != nil {
		return m.Infos
	}
	return nil
}

type FileInfo struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Size int64  `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Mode int32  `protobuf:"varint,3,opt,name=mode" json:"mode,omitempty"`
	// NOTE: int64 time containes nanoseconds from 1970-01-01T00:00:00Z
	ModTime int64  `protobuf:"varint,4,opt,name=modTime" json:"modTime,omitempty"`
	Owner   string `protobuf:"bytes,5,opt,name=owner" json:"owner,omitempty"`
	Group   string `protobuf:"bytes,6,opt,name=group" json:"group,omitempty"`
}

func (m *FileInfo) Reset()                    { *m = FileInfo{} }
func (m *FileInfo) String() string            { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()               {}
func (*FileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type FetchFileRequest struct {
	Path    string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	BufSize int32  `protobuf:"varint,2,opt,name=bufSize" json:"bufSize,omitempty"`
}

func (m *FetchFileRequest) Reset()                    { *m = FetchFileRequest{} }
func (m *FetchFileRequest) String() string            { return proto.CompactTextString(m) }
func (*FetchFileRequest) ProtoMessage()               {}
func (*FetchFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type FileChunk struct {
	Chunk []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (m *FileChunk) Reset()                    { *m = FileChunk{} }
func (m *FileChunk) String() string            { return proto.CompactTextString(m) }
func (*FileChunk) ProtoMessage()               {}
func (*FileChunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type SendFileRequest struct {
	Path  string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Chunk []byte `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (m *SendFileRequest) Reset()                    { *m = SendFileRequest{} }
func (m *SendFileRequest) String() string            { return proto.CompactTextString(m) }
func (*SendFileRequest) ProtoMessage()               {}
func (*SendFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type ChownRequest struct {
	Path  string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Owner string `protobuf:"bytes,2,opt,name=owner" json:"owner,omitempty"`
	Group string `protobuf:"bytes,3,opt,name=group" json:"group,omitempty"`
}

func (m *ChownRequest) Reset()                    { *m = ChownRequest{} }
func (m *ChownRequest) String() string            { return proto.CompactTextString(m) }
func (*ChownRequest) ProtoMessage()               {}
func (*ChownRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type ChmodRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Mode int32  `protobuf:"varint,2,opt,name=mode" json:"mode,omitempty"`
}

func (m *ChmodRequest) Reset()                    { *m = ChmodRequest{} }
func (m *ChmodRequest) String() string            { return proto.CompactTextString(m) }
func (*ChmodRequest) ProtoMessage()               {}
func (*ChmodRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type ChtimesRequest struct {
	Path  string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Atime int64  `protobuf:"varint,2,opt,name=atime" json:"atime,omitempty"`
	Mtime int64  `protobuf:"varint,3,opt,name=mtime" json:"mtime,omitempty"`
}

func (m *ChtimesRequest) Reset()                    { *m = ChtimesRequest{} }
func (m *ChtimesRequest) String() string            { return proto.CompactTextString(m) }
func (*ChtimesRequest) ProtoMessage()               {}
func (*ChtimesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type ChangeAttributesRequest struct {
	Path         string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	ChangesOwner bool   `protobuf:"varint,2,opt,name=changesOwner" json:"changesOwner,omitempty"`
	ChangesMode  bool   `protobuf:"varint,3,opt,name=changesMode" json:"changesMode,omitempty"`
	ChangesTime  bool   `protobuf:"varint,4,opt,name=changesTime" json:"changesTime,omitempty"`
	Owner        string `protobuf:"bytes,5,opt,name=owner" json:"owner,omitempty"`
	Group        string `protobuf:"bytes,6,opt,name=group" json:"group,omitempty"`
	Mode         int32  `protobuf:"varint,7,opt,name=mode" json:"mode,omitempty"`
	Atime        int64  `protobuf:"varint,8,opt,name=atime" json:"atime,omitempty"`
	Mtime        int64  `protobuf:"varint,9,opt,name=mtime" json:"mtime,omitempty"`
}

func (m *ChangeAttributesRequest) Reset()                    { *m = ChangeAttributesRequest{} }
func (m *ChangeAttributesRequest) String() string            { return proto.CompactTextString(m) }
func (*ChangeAttributesRequest) ProtoMessage()               {}
func (*ChangeAttributesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type EnsureNotExistRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *EnsureNotExistRequest) Reset()                    { *m = EnsureNotExistRequest{} }
func (m *EnsureNotExistRequest) String() string            { return proto.CompactTextString(m) }
func (*EnsureNotExistRequest) ProtoMessage()               {}
func (*EnsureNotExistRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type EnsureDirExistsRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *EnsureDirExistsRequest) Reset()                    { *m = EnsureDirExistsRequest{} }
func (m *EnsureDirExistsRequest) String() string            { return proto.CompactTextString(m) }
func (*EnsureDirExistsRequest) ProtoMessage()               {}
func (*EnsureDirExistsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func init() {
	proto.RegisterType((*StatRequest)(nil), "pb.StatRequest")
	proto.RegisterType((*ReadDirRequest)(nil), "pb.ReadDirRequest")
	proto.RegisterType((*FileInfos)(nil), "pb.FileInfos")
	proto.RegisterType((*FileInfo)(nil), "pb.FileInfo")
	proto.RegisterType((*FetchFileRequest)(nil), "pb.FetchFileRequest")
	proto.RegisterType((*FileChunk)(nil), "pb.FileChunk")
	proto.RegisterType((*SendFileRequest)(nil), "pb.SendFileRequest")
	proto.RegisterType((*Empty)(nil), "pb.Empty")
	proto.RegisterType((*ChownRequest)(nil), "pb.ChownRequest")
	proto.RegisterType((*ChmodRequest)(nil), "pb.ChmodRequest")
	proto.RegisterType((*ChtimesRequest)(nil), "pb.ChtimesRequest")
	proto.RegisterType((*ChangeAttributesRequest)(nil), "pb.ChangeAttributesRequest")
	proto.RegisterType((*EnsureNotExistRequest)(nil), "pb.EnsureNotExistRequest")
	proto.RegisterType((*EnsureDirExistsRequest)(nil), "pb.EnsureDirExistsRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RDirSync service

type RDirSyncClient interface {
	Stat(ctx context.Context, in *StatRequest, opts ...grpc.CallOption) (*FileInfo, error)
	ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (RDirSync_ReadDirClient, error)
	FetchFile(ctx context.Context, in *FetchFileRequest, opts ...grpc.CallOption) (RDirSync_FetchFileClient, error)
	SendFile(ctx context.Context, opts ...grpc.CallOption) (RDirSync_SendFileClient, error)
	Chown(ctx context.Context, in *ChownRequest, opts ...grpc.CallOption) (*Empty, error)
	Chmod(ctx context.Context, in *ChmodRequest, opts ...grpc.CallOption) (*Empty, error)
	Chtimes(ctx context.Context, in *ChtimesRequest, opts ...grpc.CallOption) (*Empty, error)
	ChangeAttributes(ctx context.Context, in *ChangeAttributesRequest, opts ...grpc.CallOption) (*Empty, error)
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

func (c *rDirSyncClient) Chown(ctx context.Context, in *ChownRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/pb.RDirSync/Chown", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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

func (c *rDirSyncClient) ChangeAttributes(ctx context.Context, in *ChangeAttributesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/pb.RDirSync/ChangeAttributes", in, out, c.cc, opts...)
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
	Chown(context.Context, *ChownRequest) (*Empty, error)
	Chmod(context.Context, *ChmodRequest) (*Empty, error)
	Chtimes(context.Context, *ChtimesRequest) (*Empty, error)
	ChangeAttributes(context.Context, *ChangeAttributesRequest) (*Empty, error)
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

func _RDirSync_Chown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDirSyncServer).Chown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RDirSync/Chown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).Chown(ctx, req.(*ChownRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _RDirSync_ChangeAttributes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeAttributesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDirSyncServer).ChangeAttributes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RDirSync/ChangeAttributes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDirSyncServer).ChangeAttributes(ctx, req.(*ChangeAttributesRequest))
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
			MethodName: "Chown",
			Handler:    _RDirSync_Chown_Handler,
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
			MethodName: "ChangeAttributes",
			Handler:    _RDirSync_ChangeAttributes_Handler,
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
	Metadata: "rdirsync.proto",
}

func init() { proto.RegisterFile("rdirsync.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 671 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x55, 0xd1, 0x6f, 0xd3, 0x3e,
	0x10, 0x5e, 0x9a, 0x65, 0x6d, 0x6f, 0x5d, 0x37, 0xf9, 0xb7, 0x1f, 0x84, 0xf2, 0x52, 0xfc, 0x00,
	0x15, 0x4c, 0x61, 0x1a, 0x08, 0x21, 0x40, 0x82, 0xad, 0xdb, 0x10, 0x0f, 0x1b, 0x23, 0xe5, 0x1f,
	0x48, 0x1a, 0x6f, 0xb1, 0xb6, 0xd8, 0xc1, 0x71, 0x18, 0xe3, 0x8d, 0x77, 0xfe, 0x63, 0x5e, 0x90,
	0xed, 0xa6, 0x8d, 0xab, 0xae, 0x13, 0x6f, 0x77, 0xdf, 0x7d, 0x77, 0x3e, 0x5f, 0xbe, 0x73, 0xa0,
	0x2b, 0x12, 0x2a, 0x8a, 0x1b, 0x36, 0x0e, 0x72, 0xc1, 0x25, 0x47, 0x8d, 0x3c, 0xc6, 0x5f, 0x60,
	0x7d, 0x24, 0x23, 0x19, 0x92, 0x6f, 0x25, 0x29, 0x24, 0x42, 0xb0, 0x9a, 0x47, 0x32, 0xf5, 0x9d,
	0xbe, 0x33, 0x68, 0x87, 0xda, 0x46, 0x01, 0xa0, 0xeb, 0x88, 0xc9, 0xe2, 0xf3, 0x35, 0x23, 0x62,
	0x9f, 0x25, 0x1f, 0x05, 0x2f, 0x73, 0xbf, 0xd1, 0x77, 0x06, 0xad, 0x70, 0x41, 0x04, 0x7f, 0x87,
	0x6e, 0x48, 0xa2, 0xe4, 0x90, 0x8a, 0x65, 0x55, 0xfb, 0xb0, 0x1e, 0xc9, 0x13, 0x5e, 0xc8, 0x21,
	0x2f, 0x99, 0xd4, 0xe5, 0xbc, 0xb0, 0x0e, 0xdd, 0x72, 0xae, 0x7b, 0xeb, 0xb9, 0xcf, 0xa1, 0x7d,
	0x4c, 0xaf, 0xc8, 0x27, 0x76, 0xce, 0x0b, 0x84, 0xc1, 0xa3, 0xca, 0xf0, 0x9d, 0xbe, 0x3b, 0x58,
	0xdf, 0xeb, 0x04, 0x79, 0x1c, 0x54, 0xd1, 0xd0, 0x84, 0xf0, 0x6f, 0x07, 0x5a, 0x15, 0xa6, 0x7a,
	0x64, 0x51, 0x46, 0xaa, 0x1e, 0x95, 0xad, 0xb0, 0x82, 0xfe, 0x24, 0xba, 0x39, 0x37, 0xd4, 0xb6,
	0xc2, 0x32, 0x9e, 0x10, 0xdd, 0x87, 0x17, 0x6a, 0x1b, 0xf9, 0xd0, 0xcc, 0x78, 0xf2, 0x95, 0x66,
	0xc4, 0x5f, 0xd5, 0xd4, 0xca, 0x45, 0xdb, 0xe0, 0x71, 0xd5, 0xa4, 0xef, 0xe9, 0xb2, 0xc6, 0x51,
	0xe8, 0x85, 0xbe, 0xcc, 0x9a, 0x41, 0xb5, 0x83, 0x3f, 0xc0, 0xd6, 0x31, 0x91, 0xe3, 0x54, 0xb5,
	0xb4, 0x6c, 0x72, 0x3e, 0x34, 0xe3, 0xf2, 0x7c, 0x54, 0x35, 0xe6, 0x85, 0x95, 0x8b, 0x1f, 0x99,
	0x09, 0x0c, 0xd3, 0x92, 0x5d, 0xaa, 0x43, 0xc6, 0xca, 0xd0, 0xb9, 0x9d, 0xd0, 0x38, 0xf8, 0x2d,
	0x6c, 0x8e, 0x08, 0x4b, 0xee, 0x3a, 0x63, 0x9a, 0xec, 0xd6, 0x93, 0x9b, 0xe0, 0x1d, 0x65, 0xb9,
	0xbc, 0xc1, 0xa7, 0xd0, 0x19, 0xa6, 0xfc, 0x9a, 0xdd, 0x51, 0xc2, 0x5c, 0xbd, 0xb1, 0xf0, 0xea,
	0x6e, 0xfd, 0xea, 0xaf, 0x54, 0xbd, 0x8c, 0x27, 0xcb, 0xea, 0x55, 0x83, 0x6f, 0xcc, 0x06, 0x8f,
	0xcf, 0xa0, 0x3b, 0x4c, 0x25, 0xcd, 0x48, 0x71, 0x47, 0x27, 0x91, 0x22, 0x4d, 0xbe, 0xa3, 0x71,
	0x14, 0x9a, 0x69, 0xd4, 0x35, 0xa8, 0x76, 0xf0, 0xaf, 0x06, 0xdc, 0x1f, 0xa6, 0x11, 0xbb, 0x20,
	0xfb, 0x52, 0x0a, 0x1a, 0x97, 0x72, 0x79, 0x6d, 0x0c, 0x9d, 0xb1, 0xa6, 0x1b, 0x31, 0x4e, 0xd6,
	0xc2, 0xc2, 0x94, 0xd4, 0x27, 0xfe, 0x49, 0xa5, 0x9c, 0x56, 0x58, 0x87, 0x6a, 0x8c, 0xa9, 0x88,
	0x66, 0x8c, 0x7f, 0x15, 0xd2, 0x74, 0x52, 0xcd, 0x9a, 0x44, 0xa7, 0x33, 0x68, 0x2d, 0x9c, 0x41,
	0xbb, 0x3e, 0x83, 0x67, 0xf0, 0xff, 0x11, 0x2b, 0x4a, 0x41, 0x4e, 0xb9, 0x3c, 0xfa, 0x41, 0x8b,
	0x65, 0xaf, 0x03, 0xde, 0x81, 0x7b, 0x86, 0x7c, 0x48, 0x85, 0x26, 0x2f, 0x1b, 0xd7, 0xde, 0x1f,
	0x17, 0x5a, 0xe1, 0x21, 0x15, 0xa3, 0x1b, 0x36, 0x46, 0x4f, 0x60, 0x55, 0xbd, 0x3d, 0x68, 0x53,
	0x2d, 0x67, 0xed, 0x15, 0xea, 0x59, 0xdb, 0x8a, 0x57, 0xd0, 0x2e, 0x34, 0x27, 0x2f, 0x0a, 0x42,
	0x2a, 0x64, 0x3f, 0x2f, 0xbd, 0x8d, 0x3a, 0xbd, 0xc0, 0x2b, 0xbb, 0x0e, 0x7a, 0x09, 0xed, 0xe9,
	0x2e, 0xa1, 0x6d, 0x1d, 0x9f, 0x5b, 0xad, 0x59, 0x96, 0x5e, 0x17, 0x9d, 0x15, 0x40, 0xab, 0x5a,
	0x0e, 0xf4, 0x9f, 0x6e, 0xca, 0x5e, 0x95, 0x5e, 0x5b, 0x81, 0x66, 0x05, 0x56, 0x06, 0x0e, 0x7a,
	0x0c, 0x9e, 0x5e, 0x03, 0xb4, 0xa5, 0xf0, 0xfa, 0x46, 0x58, 0x4c, 0xc3, 0xcb, 0x78, 0x52, 0xf1,
	0x66, 0x4a, 0xb7, 0x79, 0x4f, 0xa1, 0x39, 0x91, 0xb3, 0xb9, 0xa7, 0xad, 0x6d, 0x9b, 0xfb, 0x0e,
	0xb6, 0xe6, 0x75, 0x8a, 0x1e, 0x9a, 0xa4, 0x85, 0xea, 0xb5, 0xb3, 0xdf, 0xc0, 0xe6, 0xdc, 0x57,
	0x43, 0x3d, 0x1d, 0x5f, 0xf8, 0x29, 0xed, 0xdc, 0xd7, 0xd0, 0xb5, 0xe5, 0x81, 0x1e, 0xcc, 0x52,
	0xe7, 0x24, 0x63, 0x65, 0x1e, 0xbc, 0x87, 0x9d, 0x31, 0xcf, 0x82, 0x0b, 0x2a, 0xd3, 0x32, 0x0e,
	0x94, 0x99, 0xb2, 0xe8, 0x32, 0xca, 0x4a, 0x11, 0x4c, 0x7f, 0x4d, 0x94, 0x49, 0x22, 0x58, 0x74,
	0x15, 0xe4, 0xf1, 0xc1, 0x46, 0x25, 0x95, 0x33, 0xf5, 0xbf, 0x3a, 0x73, 0xe2, 0x35, 0xfd, 0xe3,
	0x7a, 0xf1, 0x37, 0x00, 0x00, 0xff, 0xff, 0x17, 0x07, 0x3c, 0xb0, 0xca, 0x06, 0x00, 0x00,
}
