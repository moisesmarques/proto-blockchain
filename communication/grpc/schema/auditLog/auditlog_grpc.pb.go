// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: communication/grpc/schema/auditlog.proto

package grpcAuditLog

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

// AuditLogServiceClient is the client API for AuditLogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuditLogServiceClient interface {
	// Sends a AuditLog Request
	SendAuditToAdmin(ctx context.Context, in *SendAuditToAdminRequest, opts ...grpc.CallOption) (*SendAuditToAdminResponse, error)
}

type auditLogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuditLogServiceClient(cc grpc.ClientConnInterface) AuditLogServiceClient {
	return &auditLogServiceClient{cc}
}

func (c *auditLogServiceClient) SendAuditToAdmin(ctx context.Context, in *SendAuditToAdminRequest, opts ...grpc.CallOption) (*SendAuditToAdminResponse, error) {
	out := new(SendAuditToAdminResponse)
	err := c.cc.Invoke(ctx, "/proto.AuditLogService/SendAuditToAdmin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuditLogServiceServer is the server API for AuditLogService service.
// All implementations must embed UnimplementedAuditLogServiceServer
// for forward compatibility
type AuditLogServiceServer interface {
	// Sends a AuditLog Request
	SendAuditToAdmin(context.Context, *SendAuditToAdminRequest) (*SendAuditToAdminResponse, error)
	mustEmbedUnimplementedAuditLogServiceServer()
}

// UnimplementedAuditLogServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuditLogServiceServer struct {
}

func (UnimplementedAuditLogServiceServer) SendAuditToAdmin(context.Context, *SendAuditToAdminRequest) (*SendAuditToAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendAuditToAdmin not implemented")
}
func (UnimplementedAuditLogServiceServer) mustEmbedUnimplementedAuditLogServiceServer() {}

// UnsafeAuditLogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuditLogServiceServer will
// result in compilation errors.
type UnsafeAuditLogServiceServer interface {
	mustEmbedUnimplementedAuditLogServiceServer()
}

func RegisterAuditLogServiceServer(s grpc.ServiceRegistrar, srv AuditLogServiceServer) {
	s.RegisterService(&AuditLogService_ServiceDesc, srv)
}

func _AuditLogService_SendAuditToAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendAuditToAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuditLogServiceServer).SendAuditToAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AuditLogService/SendAuditToAdmin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuditLogServiceServer).SendAuditToAdmin(ctx, req.(*SendAuditToAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuditLogService_ServiceDesc is the grpc.ServiceDesc for AuditLogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuditLogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AuditLogService",
	HandlerType: (*AuditLogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendAuditToAdmin",
			Handler:    _AuditLogService_SendAuditToAdmin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "communication/grpc/schema/auditlog.proto",
}