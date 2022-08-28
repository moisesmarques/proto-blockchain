// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: communication/grpc/schema/cpu.proto

package cpu

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

// CPUShardClient is the client API for CPUShard service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CPUShardClient interface {
	// Sends a greeting
	ProcessRequest(ctx context.Context, in *CPURequest, opts ...grpc.CallOption) (*CPUConsensusResult, error)
}

type cPUShardClient struct {
	cc grpc.ClientConnInterface
}

func NewCPUShardClient(cc grpc.ClientConnInterface) CPUShardClient {
	return &cPUShardClient{cc}
}

func (c *cPUShardClient) ProcessRequest(ctx context.Context, in *CPURequest, opts ...grpc.CallOption) (*CPUConsensusResult, error) {
	out := new(CPUConsensusResult)
	err := c.cc.Invoke(ctx, "/cpu.CPUShard/ProcessRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CPUShardServer is the server API for CPUShard service.
// All implementations must embed UnimplementedCPUShardServer
// for forward compatibility
type CPUShardServer interface {
	// Sends a greeting
	ProcessRequest(context.Context, *CPURequest) (*CPUConsensusResult, error)
	mustEmbedUnimplementedCPUShardServer()
}

// UnimplementedCPUShardServer must be embedded to have forward compatible implementations.
type UnimplementedCPUShardServer struct {
}

func (UnimplementedCPUShardServer) ProcessRequest(context.Context, *CPURequest) (*CPUConsensusResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessRequest not implemented")
}
func (UnimplementedCPUShardServer) mustEmbedUnimplementedCPUShardServer() {}

// UnsafeCPUShardServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CPUShardServer will
// result in compilation errors.
type UnsafeCPUShardServer interface {
	mustEmbedUnimplementedCPUShardServer()
}

func RegisterCPUShardServer(s grpc.ServiceRegistrar, srv CPUShardServer) {
	s.RegisterService(&CPUShard_ServiceDesc, srv)
}

func _CPUShard_ProcessRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CPURequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CPUShardServer).ProcessRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cpu.CPUShard/ProcessRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CPUShardServer).ProcessRequest(ctx, req.(*CPURequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CPUShard_ServiceDesc is the grpc.ServiceDesc for CPUShard service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CPUShard_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cpu.CPUShard",
	HandlerType: (*CPUShardServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessRequest",
			Handler:    _CPUShard_ProcessRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "communication/grpc/schema/cpu.proto",
}