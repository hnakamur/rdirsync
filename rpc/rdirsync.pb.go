// Code generated by protoc-gen-go.
// source: rdirsync.proto
// DO NOT EDIT!

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	rdirsync.proto

It has these top-level messages:
	FetchRequest
	FileChunk
	ReadDirRequest
	FileInfos
	FileInfo
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

type FetchRequest struct {
	Path    string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	BufSize int32  `protobuf:"varint,2,opt,name=bufSize" json:"bufSize,omitempty"`
}

func (m *FetchRequest) Reset()                    { *m = FetchRequest{} }
func (m *FetchRequest) String() string            { return proto.CompactTextString(m) }
func (*FetchRequest) ProtoMessage()               {}
func (*FetchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type FileChunk struct {
	Chunk []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (m *FileChunk) Reset()                    { *m = FileChunk{} }
func (m *FileChunk) String() string            { return proto.CompactTextString(m) }
func (*FileChunk) ProtoMessage()               {}
func (*FileChunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ReadDirRequest struct {
	Path        string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	AtMostCount int32  `protobuf:"varint,2,opt,name=atMostCount" json:"atMostCount,omitempty"`
}

func (m *ReadDirRequest) Reset()                    { *m = ReadDirRequest{} }
func (m *ReadDirRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadDirRequest) ProtoMessage()               {}
func (*ReadDirRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type FileInfos struct {
	Infos []*FileInfo `protobuf:"bytes,1,rep,name=infos" json:"infos,omitempty"`
}

func (m *FileInfos) Reset()                    { *m = FileInfos{} }
func (m *FileInfos) String() string            { return proto.CompactTextString(m) }
func (*FileInfos) ProtoMessage()               {}
func (*FileInfos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
	IsDir   bool   `protobuf:"varint,5,opt,name=isDir" json:"isDir,omitempty"`
}

func (m *FileInfo) Reset()                    { *m = FileInfo{} }
func (m *FileInfo) String() string            { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()               {}
func (*FileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto.RegisterType((*FetchRequest)(nil), "rpc.FetchRequest")
	proto.RegisterType((*FileChunk)(nil), "rpc.FileChunk")
	proto.RegisterType((*ReadDirRequest)(nil), "rpc.ReadDirRequest")
	proto.RegisterType((*FileInfos)(nil), "rpc.FileInfos")
	proto.RegisterType((*FileInfo)(nil), "rpc.FileInfo")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for RDirSync service

type RDirSyncClient interface {
	FetchFile(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (RDirSync_FetchFileClient, error)
	ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (RDirSync_ReadDirClient, error)
}

type rDirSyncClient struct {
	cc *grpc.ClientConn
}

func NewRDirSyncClient(cc *grpc.ClientConn) RDirSyncClient {
	return &rDirSyncClient{cc}
}

func (c *rDirSyncClient) FetchFile(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (RDirSync_FetchFileClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[0], c.cc, "/rpc.RDirSync/FetchFile", opts...)
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

func (c *rDirSyncClient) ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (RDirSync_ReadDirClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RDirSync_serviceDesc.Streams[1], c.cc, "/rpc.RDirSync/ReadDir", opts...)
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

// Server API for RDirSync service

type RDirSyncServer interface {
	FetchFile(*FetchRequest, RDirSync_FetchFileServer) error
	ReadDir(*ReadDirRequest, RDirSync_ReadDirServer) error
}

func RegisterRDirSyncServer(s *grpc.Server, srv RDirSyncServer) {
	s.RegisterService(&_RDirSync_serviceDesc, srv)
}

func _RDirSync_FetchFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FetchRequest)
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

var _RDirSync_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.RDirSync",
	HandlerType: (*RDirSyncServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FetchFile",
			Handler:       _RDirSync_FetchFile_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ReadDir",
			Handler:       _RDirSync_ReadDir_Handler,
			ServerStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("rdirsync.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x52, 0x4d, 0x6f, 0xe2, 0x30,
	0x14, 0x5c, 0x6f, 0xc8, 0x02, 0x8f, 0x0f, 0x69, 0xdd, 0x1e, 0xac, 0x9e, 0xd2, 0x70, 0xc9, 0x29,
	0x42, 0xa9, 0xd4, 0x53, 0x4f, 0x80, 0x90, 0x7a, 0xa8, 0x84, 0x4c, 0xff, 0x40, 0xe2, 0x98, 0xc6,
	0xa2, 0xb1, 0x53, 0xc7, 0xa9, 0x44, 0x7f, 0x7d, 0x65, 0x07, 0x23, 0x7a, 0xe9, 0x6d, 0x66, 0x5e,
	0x9e, 0x66, 0xde, 0xc4, 0x30, 0xd7, 0xa5, 0xd0, 0xed, 0x49, 0xb2, 0xb4, 0xd1, 0xca, 0x28, 0x1c,
	0xe8, 0x86, 0xc5, 0x4f, 0x30, 0xdd, 0x72, 0xc3, 0x2a, 0xca, 0x3f, 0x3a, 0xde, 0x1a, 0x8c, 0x61,
	0xd0, 0xe4, 0xa6, 0x22, 0x28, 0x42, 0xc9, 0x98, 0x3a, 0x8c, 0x09, 0x0c, 0x8b, 0xee, 0xb0, 0x17,
	0x5f, 0x9c, 0xfc, 0x8d, 0x50, 0x12, 0x52, 0x4f, 0xe3, 0x7b, 0x18, 0x6f, 0xc5, 0x3b, 0x5f, 0x57,
	0x9d, 0x3c, 0xe2, 0x5b, 0x08, 0x99, 0x05, 0x6e, 0x77, 0x4a, 0x7b, 0x12, 0x6f, 0x61, 0x4e, 0x79,
	0x5e, 0x6e, 0x84, 0xfe, 0xcd, 0x22, 0x82, 0x49, 0x6e, 0x5e, 0x54, 0x6b, 0xd6, 0xaa, 0x93, 0xe6,
	0x6c, 0x73, 0x2d, 0xc5, 0xcb, 0xde, 0xea, 0x59, 0x1e, 0x54, 0x8b, 0x17, 0x10, 0x0a, 0x0b, 0x08,
	0x8a, 0x82, 0x64, 0x92, 0xcd, 0x52, 0xdd, 0xb0, 0xd4, 0x8f, 0x69, 0x3f, 0x8b, 0x3f, 0x61, 0xe4,
	0x25, 0xeb, 0x29, 0xf3, 0x9a, 0x7b, 0x4f, 0x8b, 0xad, 0xd6, 0xfa, 0x9b, 0x02, 0xea, 0xb0, 0xd5,
	0x6a, 0x55, 0x72, 0x12, 0xb8, 0x00, 0x0e, 0xdb, 0xf3, 0x6b, 0x55, 0xbe, 0x8a, 0x9a, 0x93, 0x81,
	0xfb, 0xd4, 0x53, 0x7b, 0xb1, 0x68, 0x37, 0x42, 0x93, 0x30, 0x42, 0xc9, 0x88, 0xf6, 0x24, 0xd3,
	0x30, 0xa2, 0x1b, 0xa1, 0xf7, 0x27, 0xc9, 0x70, 0x06, 0x63, 0x57, 0xaf, 0x0d, 0x82, 0xff, 0xf7,
	0x31, 0xaf, 0xea, 0xbe, 0x9b, 0x5f, 0x92, 0xbb, 0x0e, 0xe3, 0x3f, 0x4b, 0x84, 0x33, 0x18, 0x9e,
	0x1b, 0xc3, 0x37, 0x6e, 0xfc, 0xb3, 0xbf, 0xab, 0x1d, 0x57, 0x86, 0xdd, 0x59, 0x3d, 0xc2, 0x42,
	0xe9, 0xb7, 0xb4, 0x10, 0xa6, 0xe8, 0xd8, 0x91, 0x9b, 0xb4, 0x92, 0xf9, 0x31, 0xaf, 0x3b, 0x9d,
	0x5e, 0x7e, 0xba, 0x6e, 0xd8, 0x6a, 0xe6, 0x83, 0xed, 0xec, 0x0b, 0xd8, 0xa1, 0xe2, 0x9f, 0x7b,
	0x0a, 0x0f, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6e, 0xd5, 0x07, 0xd7, 0x1c, 0x02, 0x00, 0x00,
}