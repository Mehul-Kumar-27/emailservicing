// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: logs.proto

package logs

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	LogServeice_WriteLog_FullMethodName = "/LogServeice/WriteLog"
)

// LogServeiceClient is the client API for LogServeice service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogServeiceClient interface {
	WriteLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error)
}

type logServeiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogServeiceClient(cc grpc.ClientConnInterface) LogServeiceClient {
	return &logServeiceClient{cc}
}

func (c *logServeiceClient) WriteLog(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogResponse, error) {
	out := new(LogResponse)
	err := c.cc.Invoke(ctx, LogServeice_WriteLog_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogServeiceServer is the server API for LogServeice service.
// All implementations must embed UnimplementedLogServeiceServer
// for forward compatibility
type LogServeiceServer interface {
	WriteLog(context.Context, *LogRequest) (*LogResponse, error)
	mustEmbedUnimplementedLogServeiceServer()
}

// UnimplementedLogServeiceServer must be embedded to have forward compatible implementations.
type UnimplementedLogServeiceServer struct {
}

func (UnimplementedLogServeiceServer) WriteLog(context.Context, *LogRequest) (*LogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteLog not implemented")
}
func (UnimplementedLogServeiceServer) mustEmbedUnimplementedLogServeiceServer() {}

// UnsafeLogServeiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogServeiceServer will
// result in compilation errors.
type UnsafeLogServeiceServer interface {
	mustEmbedUnimplementedLogServeiceServer()
}

func RegisterLogServeiceServer(s grpc.ServiceRegistrar, srv LogServeiceServer) {
	s.RegisterService(&LogServeice_ServiceDesc, srv)
}

func _LogServeice_WriteLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServeiceServer).WriteLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogServeice_WriteLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServeiceServer).WriteLog(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LogServeice_ServiceDesc is the grpc.ServiceDesc for LogServeice service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogServeice_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "LogServeice",
	HandlerType: (*LogServeiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WriteLog",
			Handler:    _LogServeice_WriteLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logs.proto",
}