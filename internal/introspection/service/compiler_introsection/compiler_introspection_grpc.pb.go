// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: internal/introspection/service/proto/compiler_introspection.proto

package compiler_introsection

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
	CompilerIntrospection_Helo_FullMethodName         = "/compilerintrospection.CompilerIntrospection/Helo"
	CompilerIntrospection_EventStream_FullMethodName  = "/compilerintrospection.CompilerIntrospection/EventStream"
	CompilerIntrospection_AbortSession_FullMethodName = "/compilerintrospection.CompilerIntrospection/AbortSession"
)

// CompilerIntrospectionClient is the client API for CompilerIntrospection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompilerIntrospectionClient interface {
	Helo(ctx context.Context, in *HeloRequest, opts ...grpc.CallOption) (*HeloResponse, error)
	EventStream(ctx context.Context, in *EventStreamRequest, opts ...grpc.CallOption) (CompilerIntrospection_EventStreamClient, error)
	AbortSession(ctx context.Context, in *AbortSessionRequest, opts ...grpc.CallOption) (*AbortSessionResponse, error)
}

type compilerIntrospectionClient struct {
	cc grpc.ClientConnInterface
}

func NewCompilerIntrospectionClient(cc grpc.ClientConnInterface) CompilerIntrospectionClient {
	return &compilerIntrospectionClient{cc}
}

func (c *compilerIntrospectionClient) Helo(ctx context.Context, in *HeloRequest, opts ...grpc.CallOption) (*HeloResponse, error) {
	out := new(HeloResponse)
	err := c.cc.Invoke(ctx, CompilerIntrospection_Helo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compilerIntrospectionClient) EventStream(ctx context.Context, in *EventStreamRequest, opts ...grpc.CallOption) (CompilerIntrospection_EventStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &CompilerIntrospection_ServiceDesc.Streams[0], CompilerIntrospection_EventStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &compilerIntrospectionEventStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CompilerIntrospection_EventStreamClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type compilerIntrospectionEventStreamClient struct {
	grpc.ClientStream
}

func (x *compilerIntrospectionEventStreamClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *compilerIntrospectionClient) AbortSession(ctx context.Context, in *AbortSessionRequest, opts ...grpc.CallOption) (*AbortSessionResponse, error) {
	out := new(AbortSessionResponse)
	err := c.cc.Invoke(ctx, CompilerIntrospection_AbortSession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompilerIntrospectionServer is the server API for CompilerIntrospection service.
// All implementations must embed UnimplementedCompilerIntrospectionServer
// for forward compatibility
type CompilerIntrospectionServer interface {
	Helo(context.Context, *HeloRequest) (*HeloResponse, error)
	EventStream(*EventStreamRequest, CompilerIntrospection_EventStreamServer) error
	AbortSession(context.Context, *AbortSessionRequest) (*AbortSessionResponse, error)
	mustEmbedUnimplementedCompilerIntrospectionServer()
}

// UnimplementedCompilerIntrospectionServer must be embedded to have forward compatible implementations.
type UnimplementedCompilerIntrospectionServer struct {
}

func (UnimplementedCompilerIntrospectionServer) Helo(context.Context, *HeloRequest) (*HeloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Helo not implemented")
}
func (UnimplementedCompilerIntrospectionServer) EventStream(*EventStreamRequest, CompilerIntrospection_EventStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method EventStream not implemented")
}
func (UnimplementedCompilerIntrospectionServer) AbortSession(context.Context, *AbortSessionRequest) (*AbortSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AbortSession not implemented")
}
func (UnimplementedCompilerIntrospectionServer) mustEmbedUnimplementedCompilerIntrospectionServer() {}

// UnsafeCompilerIntrospectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompilerIntrospectionServer will
// result in compilation errors.
type UnsafeCompilerIntrospectionServer interface {
	mustEmbedUnimplementedCompilerIntrospectionServer()
}

func RegisterCompilerIntrospectionServer(s grpc.ServiceRegistrar, srv CompilerIntrospectionServer) {
	s.RegisterService(&CompilerIntrospection_ServiceDesc, srv)
}

func _CompilerIntrospection_Helo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompilerIntrospectionServer).Helo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompilerIntrospection_Helo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompilerIntrospectionServer).Helo(ctx, req.(*HeloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CompilerIntrospection_EventStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EventStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CompilerIntrospectionServer).EventStream(m, &compilerIntrospectionEventStreamServer{stream})
}

type CompilerIntrospection_EventStreamServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type compilerIntrospectionEventStreamServer struct {
	grpc.ServerStream
}

func (x *compilerIntrospectionEventStreamServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func _CompilerIntrospection_AbortSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AbortSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompilerIntrospectionServer).AbortSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompilerIntrospection_AbortSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompilerIntrospectionServer).AbortSession(ctx, req.(*AbortSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CompilerIntrospection_ServiceDesc is the grpc.ServiceDesc for CompilerIntrospection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CompilerIntrospection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "compilerintrospection.CompilerIntrospection",
	HandlerType: (*CompilerIntrospectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Helo",
			Handler:    _CompilerIntrospection_Helo_Handler,
		},
		{
			MethodName: "AbortSession",
			Handler:    _CompilerIntrospection_AbortSession_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EventStream",
			Handler:       _CompilerIntrospection_EventStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/introspection/service/proto/compiler_introspection.proto",
}
