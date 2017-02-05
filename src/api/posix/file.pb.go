// Code generated by protoc-gen-go.
// source: file.proto
// DO NOT EDIT!

/*
Package posix is a generated protocol buffer package.

It is generated from these files:
	file.proto

It has these top-level messages:
	Info
	Filter
	List
	Chunk
	LocalOperation
	RemoteOperation
*/
package posix

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

type Info_State int32

const (
	Info_OK    Info_State = 0
	Info_ERROR Info_State = 1
)

var Info_State_name = map[int32]string{
	0: "OK",
	1: "ERROR",
}
var Info_State_value = map[string]int32{
	"OK":    0,
	"ERROR": 1,
}

func (x Info_State) String() string {
	return proto.EnumName(Info_State_name, int32(x))
}
func (Info_State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type List_State int32

const (
	List_OK    List_State = 0
	List_RETRY List_State = 1
	List_ERROR List_State = 2
)

var List_State_name = map[int32]string{
	0: "OK",
	1: "RETRY",
	2: "ERROR",
}
var List_State_value = map[string]int32{
	"OK":    0,
	"RETRY": 1,
	"ERROR": 2,
}

func (x List_State) String() string {
	return proto.EnumName(List_State_name, int32(x))
}
func (List_State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

type LocalOperation_Type int32

const (
	LocalOperation_DELETE LocalOperation_Type = 0
	LocalOperation_MOVE   LocalOperation_Type = 1
	LocalOperation_COPY   LocalOperation_Type = 2
)

var LocalOperation_Type_name = map[int32]string{
	0: "DELETE",
	1: "MOVE",
	2: "COPY",
}
var LocalOperation_Type_value = map[string]int32{
	"DELETE": 0,
	"MOVE":   1,
	"COPY":   2,
}

func (x LocalOperation_Type) String() string {
	return proto.EnumName(LocalOperation_Type_name, int32(x))
}
func (LocalOperation_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

type LocalOperation_State int32

const (
	LocalOperation_OK    LocalOperation_State = 0
	LocalOperation_ERROR LocalOperation_State = 1
)

var LocalOperation_State_name = map[int32]string{
	0: "OK",
	1: "ERROR",
}
var LocalOperation_State_value = map[string]int32{
	"OK":    0,
	"ERROR": 1,
}

func (x LocalOperation_State) String() string {
	return proto.EnumName(LocalOperation_State_name, int32(x))
}
func (LocalOperation_State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 1} }

type RemoteOperation_Type int32

const (
	RemoteOperation_COPY_FROM RemoteOperation_Type = 0
)

var RemoteOperation_Type_name = map[int32]string{
	0: "COPY_FROM",
}
var RemoteOperation_Type_value = map[string]int32{
	"COPY_FROM": 0,
}

func (x RemoteOperation_Type) String() string {
	return proto.EnumName(RemoteOperation_Type_name, int32(x))
}
func (RemoteOperation_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

type RemoteOperation_State int32

const (
	RemoteOperation_OK    RemoteOperation_State = 0
	RemoteOperation_ERROR RemoteOperation_State = 1
)

var RemoteOperation_State_name = map[int32]string{
	0: "OK",
	1: "ERROR",
}
var RemoteOperation_State_value = map[string]int32{
	"OK":    0,
	"ERROR": 1,
}

func (x RemoteOperation_State) String() string {
	return proto.EnumName(RemoteOperation_State_name, int32(x))
}
func (RemoteOperation_State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 1} }

type Info struct {
	Path  string     `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Size  int64      `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Md5   string     `protobuf:"bytes,3,opt,name=md5" json:"md5,omitempty"`
	State Info_State `protobuf:"varint,4,opt,name=state,enum=posix.Info_State" json:"state,omitempty"`
	Error string     `protobuf:"bytes,5,opt,name=error" json:"error,omitempty"`
}

func (m *Info) Reset()                    { *m = Info{} }
func (m *Info) String() string            { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()               {}
func (*Info) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Info) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Info) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *Info) GetMd5() string {
	if m != nil {
		return m.Md5
	}
	return ""
}

func (m *Info) GetState() Info_State {
	if m != nil {
		return m.State
	}
	return Info_OK
}

func (m *Info) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type Filter struct {
	PathMatch string `protobuf:"bytes,1,opt,name=path_match,json=pathMatch" json:"path_match,omitempty"`
}

func (m *Filter) Reset()                    { *m = Filter{} }
func (m *Filter) String() string            { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()               {}
func (*Filter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Filter) GetPathMatch() string {
	if m != nil {
		return m.PathMatch
	}
	return ""
}

type List struct {
	Filter      *Filter    `protobuf:"bytes,1,opt,name=filter" json:"filter,omitempty"`
	SnaphotTime int64      `protobuf:"varint,2,opt,name=snaphot_time,json=snaphotTime" json:"snaphot_time,omitempty"`
	Files       []*Info    `protobuf:"bytes,3,rep,name=files" json:"files,omitempty"`
	State       List_State `protobuf:"varint,4,opt,name=state,enum=posix.List_State" json:"state,omitempty"`
	Error       string     `protobuf:"bytes,5,opt,name=error" json:"error,omitempty"`
}

func (m *List) Reset()                    { *m = List{} }
func (m *List) String() string            { return proto.CompactTextString(m) }
func (*List) ProtoMessage()               {}
func (*List) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *List) GetFilter() *Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *List) GetSnaphotTime() int64 {
	if m != nil {
		return m.SnaphotTime
	}
	return 0
}

func (m *List) GetFiles() []*Info {
	if m != nil {
		return m.Files
	}
	return nil
}

func (m *List) GetState() List_State {
	if m != nil {
		return m.State
	}
	return List_OK
}

func (m *List) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type Chunk struct {
	File   *Info  `protobuf:"bytes,1,opt,name=file" json:"file,omitempty"`
	Offset int64  `protobuf:"varint,2,opt,name=offset" json:"offset,omitempty"`
	Data   []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Chunk) Reset()                    { *m = Chunk{} }
func (m *Chunk) String() string            { return proto.CompactTextString(m) }
func (*Chunk) ProtoMessage()               {}
func (*Chunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Chunk) GetFile() *Info {
	if m != nil {
		return m.File
	}
	return nil
}

func (m *Chunk) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *Chunk) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type LocalOperation struct {
	Type    LocalOperation_Type  `protobuf:"varint,1,opt,name=type,enum=posix.LocalOperation_Type" json:"type,omitempty"`
	File    *Info                `protobuf:"bytes,2,opt,name=file" json:"file,omitempty"`
	DstFile *Info                `protobuf:"bytes,3,opt,name=dst_file,json=dstFile" json:"dst_file,omitempty"`
	State   LocalOperation_State `protobuf:"varint,4,opt,name=state,enum=posix.LocalOperation_State" json:"state,omitempty"`
	Error   string               `protobuf:"bytes,5,opt,name=error" json:"error,omitempty"`
}

func (m *LocalOperation) Reset()                    { *m = LocalOperation{} }
func (m *LocalOperation) String() string            { return proto.CompactTextString(m) }
func (*LocalOperation) ProtoMessage()               {}
func (*LocalOperation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *LocalOperation) GetType() LocalOperation_Type {
	if m != nil {
		return m.Type
	}
	return LocalOperation_DELETE
}

func (m *LocalOperation) GetFile() *Info {
	if m != nil {
		return m.File
	}
	return nil
}

func (m *LocalOperation) GetDstFile() *Info {
	if m != nil {
		return m.DstFile
	}
	return nil
}

func (m *LocalOperation) GetState() LocalOperation_State {
	if m != nil {
		return m.State
	}
	return LocalOperation_OK
}

func (m *LocalOperation) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type RemoteOperation struct {
	Type         RemoteOperation_Type  `protobuf:"varint,1,opt,name=type,enum=posix.RemoteOperation_Type" json:"type,omitempty"`
	ToFile       *Info                 `protobuf:"bytes,2,opt,name=to_file,json=toFile" json:"to_file,omitempty"`
	RemoteServer string                `protobuf:"bytes,3,opt,name=remote_server,json=remoteServer" json:"remote_server,omitempty"`
	RemotePort   int64                 `protobuf:"varint,4,opt,name=remote_port,json=remotePort" json:"remote_port,omitempty"`
	RemoteFile   *Info                 `protobuf:"bytes,5,opt,name=remote_file,json=remoteFile" json:"remote_file,omitempty"`
	State        RemoteOperation_State `protobuf:"varint,6,opt,name=state,enum=posix.RemoteOperation_State" json:"state,omitempty"`
	Error        string                `protobuf:"bytes,7,opt,name=error" json:"error,omitempty"`
}

func (m *RemoteOperation) Reset()                    { *m = RemoteOperation{} }
func (m *RemoteOperation) String() string            { return proto.CompactTextString(m) }
func (*RemoteOperation) ProtoMessage()               {}
func (*RemoteOperation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RemoteOperation) GetType() RemoteOperation_Type {
	if m != nil {
		return m.Type
	}
	return RemoteOperation_COPY_FROM
}

func (m *RemoteOperation) GetToFile() *Info {
	if m != nil {
		return m.ToFile
	}
	return nil
}

func (m *RemoteOperation) GetRemoteServer() string {
	if m != nil {
		return m.RemoteServer
	}
	return ""
}

func (m *RemoteOperation) GetRemotePort() int64 {
	if m != nil {
		return m.RemotePort
	}
	return 0
}

func (m *RemoteOperation) GetRemoteFile() *Info {
	if m != nil {
		return m.RemoteFile
	}
	return nil
}

func (m *RemoteOperation) GetState() RemoteOperation_State {
	if m != nil {
		return m.State
	}
	return RemoteOperation_OK
}

func (m *RemoteOperation) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*Info)(nil), "posix.Info")
	proto.RegisterType((*Filter)(nil), "posix.Filter")
	proto.RegisterType((*List)(nil), "posix.List")
	proto.RegisterType((*Chunk)(nil), "posix.Chunk")
	proto.RegisterType((*LocalOperation)(nil), "posix.LocalOperation")
	proto.RegisterType((*RemoteOperation)(nil), "posix.RemoteOperation")
	proto.RegisterEnum("posix.Info_State", Info_State_name, Info_State_value)
	proto.RegisterEnum("posix.List_State", List_State_name, List_State_value)
	proto.RegisterEnum("posix.LocalOperation_Type", LocalOperation_Type_name, LocalOperation_Type_value)
	proto.RegisterEnum("posix.LocalOperation_State", LocalOperation_State_name, LocalOperation_State_value)
	proto.RegisterEnum("posix.RemoteOperation_Type", RemoteOperation_Type_name, RemoteOperation_Type_value)
	proto.RegisterEnum("posix.RemoteOperation_State", RemoteOperation_State_name, RemoteOperation_State_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for File service

type FileClient interface {
	Find(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*List, error)
	Receive(ctx context.Context, opts ...grpc.CallOption) (File_ReceiveClient, error)
	Stream(ctx context.Context, in *Info, opts ...grpc.CallOption) (File_StreamClient, error)
	LocalOps(ctx context.Context, in *LocalOperation, opts ...grpc.CallOption) (*LocalOperation, error)
	RemoteOps(ctx context.Context, in *RemoteOperation, opts ...grpc.CallOption) (*RemoteOperation, error)
}

type fileClient struct {
	cc *grpc.ClientConn
}

func NewFileClient(cc *grpc.ClientConn) FileClient {
	return &fileClient{cc}
}

func (c *fileClient) Find(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*List, error) {
	out := new(List)
	err := grpc.Invoke(ctx, "/posix.File/Find", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) Receive(ctx context.Context, opts ...grpc.CallOption) (File_ReceiveClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_File_serviceDesc.Streams[0], c.cc, "/posix.File/Receive", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileReceiveClient{stream}
	return x, nil
}

type File_ReceiveClient interface {
	Send(*Chunk) error
	CloseAndRecv() (*Info, error)
	grpc.ClientStream
}

type fileReceiveClient struct {
	grpc.ClientStream
}

func (x *fileReceiveClient) Send(m *Chunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileReceiveClient) CloseAndRecv() (*Info, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Info)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileClient) Stream(ctx context.Context, in *Info, opts ...grpc.CallOption) (File_StreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_File_serviceDesc.Streams[1], c.cc, "/posix.File/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type File_StreamClient interface {
	Recv() (*Chunk, error)
	grpc.ClientStream
}

type fileStreamClient struct {
	grpc.ClientStream
}

func (x *fileStreamClient) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileClient) LocalOps(ctx context.Context, in *LocalOperation, opts ...grpc.CallOption) (*LocalOperation, error) {
	out := new(LocalOperation)
	err := grpc.Invoke(ctx, "/posix.File/LocalOps", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileClient) RemoteOps(ctx context.Context, in *RemoteOperation, opts ...grpc.CallOption) (*RemoteOperation, error) {
	out := new(RemoteOperation)
	err := grpc.Invoke(ctx, "/posix.File/RemoteOps", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for File service

type FileServer interface {
	Find(context.Context, *Filter) (*List, error)
	Receive(File_ReceiveServer) error
	Stream(*Info, File_StreamServer) error
	LocalOps(context.Context, *LocalOperation) (*LocalOperation, error)
	RemoteOps(context.Context, *RemoteOperation) (*RemoteOperation, error)
}

func RegisterFileServer(s *grpc.Server, srv FileServer) {
	s.RegisterService(&_File_serviceDesc, srv)
}

func _File_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/posix.File/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).Find(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_Receive_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServer).Receive(&fileReceiveServer{stream})
}

type File_ReceiveServer interface {
	SendAndClose(*Info) error
	Recv() (*Chunk, error)
	grpc.ServerStream
}

type fileReceiveServer struct {
	grpc.ServerStream
}

func (x *fileReceiveServer) SendAndClose(m *Info) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileReceiveServer) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _File_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Info)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServer).Stream(m, &fileStreamServer{stream})
}

type File_StreamServer interface {
	Send(*Chunk) error
	grpc.ServerStream
}

type fileStreamServer struct {
	grpc.ServerStream
}

func (x *fileStreamServer) Send(m *Chunk) error {
	return x.ServerStream.SendMsg(m)
}

func _File_LocalOps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocalOperation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).LocalOps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/posix.File/LocalOps",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).LocalOps(ctx, req.(*LocalOperation))
	}
	return interceptor(ctx, in, info, handler)
}

func _File_RemoteOps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoteOperation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).RemoteOps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/posix.File/RemoteOps",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).RemoteOps(ctx, req.(*RemoteOperation))
	}
	return interceptor(ctx, in, info, handler)
}

var _File_serviceDesc = grpc.ServiceDesc{
	ServiceName: "posix.File",
	HandlerType: (*FileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Find",
			Handler:    _File_Find_Handler,
		},
		{
			MethodName: "LocalOps",
			Handler:    _File_LocalOps_Handler,
		},
		{
			MethodName: "RemoteOps",
			Handler:    _File_RemoteOps_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Receive",
			Handler:       _File_Receive_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Stream",
			Handler:       _File_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "file.proto",
}

func init() { proto.RegisterFile("file.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 619 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x54, 0xcd, 0x6e, 0xd3, 0x4c,
	0x14, 0xf5, 0x7f, 0x92, 0x9b, 0x9f, 0x2f, 0xdf, 0x88, 0x56, 0x91, 0x01, 0x01, 0xa6, 0x3f, 0x59,
	0x20, 0x03, 0x41, 0x6c, 0x90, 0x58, 0x15, 0x47, 0x42, 0x24, 0x72, 0x35, 0x89, 0x90, 0xba, 0xb2,
	0x4c, 0x32, 0x51, 0x2c, 0x92, 0x8c, 0x65, 0x0f, 0x15, 0xf0, 0x20, 0x6c, 0x78, 0xb2, 0xbe, 0x0a,
	0x2b, 0x66, 0xc6, 0x76, 0x6b, 0x53, 0xa7, 0xea, 0x2a, 0xe3, 0x73, 0x8f, 0xef, 0x3d, 0xe7, 0x78,
	0x6e, 0x00, 0x56, 0xd1, 0x86, 0xb8, 0x71, 0x42, 0x19, 0x45, 0x66, 0x4c, 0xd3, 0xe8, 0xbb, 0xf3,
	0x5b, 0x05, 0xe3, 0xe3, 0x6e, 0x45, 0x11, 0x02, 0x23, 0x0e, 0xd9, 0x7a, 0xa0, 0x3e, 0x55, 0x87,
	0x2d, 0x2c, 0xcf, 0x02, 0x4b, 0xa3, 0x9f, 0x64, 0xa0, 0x71, 0x4c, 0xc7, 0xf2, 0x8c, 0xfa, 0xa0,
	0x6f, 0x97, 0x6f, 0x07, 0xba, 0xa4, 0x89, 0x23, 0x3a, 0x05, 0x33, 0x65, 0x21, 0x23, 0x03, 0x83,
	0x63, 0xbd, 0xd1, 0xff, 0xae, 0xec, 0xec, 0x8a, 0xae, 0xee, 0x4c, 0x14, 0x70, 0x56, 0x47, 0x0f,
	0xc0, 0x24, 0x49, 0x42, 0x93, 0x81, 0x29, 0x5f, 0xce, 0x1e, 0x1c, 0x1b, 0x4c, 0xc9, 0x42, 0x16,
	0x68, 0xfe, 0xa7, 0xbe, 0x82, 0x5a, 0x60, 0x7a, 0x18, 0xfb, 0xb8, 0xaf, 0x3a, 0xa7, 0x60, 0x8d,
	0xa3, 0x0d, 0x23, 0x09, 0x7a, 0x0c, 0x20, 0x24, 0x05, 0xdb, 0x90, 0x2d, 0x0a, 0x91, 0x2d, 0x81,
	0x4c, 0x05, 0xe0, 0x5c, 0x71, 0x1b, 0x93, 0x28, 0x65, 0xe8, 0x18, 0xac, 0x95, 0x7c, 0x43, 0x72,
	0xda, 0xa3, 0x6e, 0xae, 0x26, 0x6b, 0x83, 0xf3, 0x22, 0x7a, 0x06, 0x9d, 0x74, 0x17, 0xc6, 0x6b,
	0xca, 0x02, 0x16, 0x6d, 0x0b, 0x87, 0xed, 0x1c, 0x9b, 0x73, 0x88, 0x53, 0x4c, 0x11, 0x57, 0xca,
	0xad, 0xea, 0xbc, 0x51, 0xbb, 0x64, 0x0b, 0x67, 0x95, 0x7d, 0xce, 0x85, 0x90, 0xfb, 0x38, 0x3f,
	0xae, 0x71, 0x8e, 0xbd, 0x39, 0xbe, 0xe8, 0xab, 0x37, 0x21, 0x68, 0xce, 0x1c, 0xcc, 0xb3, 0xf5,
	0xb7, 0xdd, 0x57, 0xf4, 0x04, 0x0c, 0x31, 0x37, 0x77, 0x56, 0x11, 0x24, 0x0b, 0xe8, 0x10, 0x2c,
	0xba, 0x5a, 0xa5, 0x84, 0xe5, 0x7e, 0xf2, 0x27, 0xf1, 0x1d, 0x97, 0x21, 0x0b, 0xe5, 0x47, 0xeb,
	0x60, 0x79, 0x76, 0x7e, 0x69, 0xd0, 0x9b, 0xd0, 0x45, 0xb8, 0xf1, 0x63, 0x92, 0x84, 0x2c, 0xa2,
	0x3b, 0xe4, 0x82, 0xc1, 0x7e, 0xc4, 0x59, 0xff, 0xde, 0xc8, 0x2e, 0xdc, 0x54, 0x48, 0xee, 0x9c,
	0x33, 0xb0, 0xe4, 0x5d, 0xeb, 0xd1, 0xf6, 0xe9, 0x39, 0x81, 0xe6, 0x32, 0x65, 0x81, 0x24, 0xe9,
	0xb7, 0x49, 0x0d, 0x5e, 0x1c, 0x0b, 0xde, 0xeb, 0x6a, 0x8e, 0x0f, 0xeb, 0x27, 0xdf, 0x23, 0xd1,
	0x13, 0x30, 0x84, 0x3e, 0x04, 0x60, 0x7d, 0xf0, 0x26, 0xde, 0xdc, 0xe3, 0xa1, 0x36, 0xc1, 0x98,
	0xfa, 0x9f, 0x3d, 0x9e, 0x29, 0x3f, 0x9d, 0xf9, 0xe7, 0x17, 0x3c, 0xd2, 0xbb, 0xee, 0xdc, 0x95,
	0x06, 0xff, 0x61, 0xb2, 0xa5, 0x8c, 0xdc, 0x24, 0xf3, 0xb2, 0x92, 0x4c, 0xa1, 0xef, 0x1f, 0x56,
	0x39, 0x9a, 0x23, 0x68, 0x30, 0x1a, 0xec, 0x4b, 0xc7, 0x62, 0x54, 0xfa, 0x7e, 0x0e, 0xdd, 0x44,
	0xf6, 0x08, 0x52, 0x92, 0x5c, 0xf2, 0x3b, 0x9b, 0x6d, 0x55, 0x27, 0x03, 0x67, 0x12, 0xe3, 0x29,
	0xb7, 0x73, 0x52, 0x4c, 0x13, 0x26, 0x23, 0xd2, 0x31, 0x64, 0xd0, 0x39, 0x47, 0xd0, 0x8b, 0x6b,
	0x82, 0x9c, 0x67, 0xde, 0x9e, 0x97, 0xb3, 0xe5, 0xcc, 0x51, 0x91, 0xb5, 0x25, 0xbd, 0x3c, 0xda,
	0xe3, 0xa5, 0x3e, 0xec, 0x46, 0x39, 0xec, 0x83, 0x3c, 0xec, 0x2e, 0xb4, 0x44, 0xac, 0xc1, 0x18,
	0xfb, 0xd3, 0xbe, 0x72, 0x57, 0xb6, 0xa3, 0x3f, 0x7c, 0x4d, 0xa5, 0x8a, 0x23, 0xf1, 0xbb, 0x5b,
	0xa2, 0xea, 0x7a, 0xda, 0xed, 0xd2, 0x06, 0x39, 0x0a, 0x1a, 0x42, 0x03, 0x93, 0x05, 0x89, 0x2e,
	0x09, 0xea, 0xe4, 0x15, 0xb9, 0x09, 0x76, 0xd9, 0x9d, 0xa3, 0x0c, 0x55, 0xbe, 0x89, 0xd6, 0x8c,
	0x25, 0x24, 0xdc, 0xa2, 0x72, 0xc9, 0xae, 0xbc, 0xe5, 0x28, 0xaf, 0x54, 0xf4, 0x0e, 0x9a, 0xf9,
	0xb5, 0x4a, 0xd1, 0x41, 0xed, 0x3d, 0xb3, 0xeb, 0x61, 0x2e, 0xe7, 0x3d, 0xb4, 0x8a, 0x98, 0x52,
	0x74, 0x58, 0x1f, 0x9c, 0xbd, 0x07, 0x77, 0x94, 0x2f, 0x96, 0xfc, 0xe3, 0x7d, 0xf3, 0x37, 0x00,
	0x00, 0xff, 0xff, 0x08, 0xf7, 0xb0, 0x6e, 0x86, 0x05, 0x00, 0x00,
}
