// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: delegatepb/delegate.proto

package delegatepb

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
	Delegate_GetDelegates_FullMethodName = "/delegatepb.Delegate/GetDelegates"
)

// DelegateClient is the client API for Delegate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DelegateClient interface {
	GetDelegates(ctx context.Context, in *GetDelegatesRequest, opts ...grpc.CallOption) (*GetDelegatesResponse, error)
}

type delegateClient struct {
	cc grpc.ClientConnInterface
}

func NewDelegateClient(cc grpc.ClientConnInterface) DelegateClient {
	return &delegateClient{cc}
}

func (c *delegateClient) GetDelegates(ctx context.Context, in *GetDelegatesRequest, opts ...grpc.CallOption) (*GetDelegatesResponse, error) {
	out := new(GetDelegatesResponse)
	err := c.cc.Invoke(ctx, Delegate_GetDelegates_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DelegateServer is the server API for Delegate service.
// All implementations must embed UnimplementedDelegateServer
// for forward compatibility
type DelegateServer interface {
	GetDelegates(context.Context, *GetDelegatesRequest) (*GetDelegatesResponse, error)
	mustEmbedUnimplementedDelegateServer()
}

// UnimplementedDelegateServer must be embedded to have forward compatible implementations.
type UnimplementedDelegateServer struct {
}

func (UnimplementedDelegateServer) GetDelegates(context.Context, *GetDelegatesRequest) (*GetDelegatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDelegates not implemented")
}
func (UnimplementedDelegateServer) mustEmbedUnimplementedDelegateServer() {}

// UnsafeDelegateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DelegateServer will
// result in compilation errors.
type UnsafeDelegateServer interface {
	mustEmbedUnimplementedDelegateServer()
}

func RegisterDelegateServer(s grpc.ServiceRegistrar, srv DelegateServer) {
	s.RegisterService(&Delegate_ServiceDesc, srv)
}

func _Delegate_GetDelegates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDelegatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DelegateServer).GetDelegates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Delegate_GetDelegates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DelegateServer).GetDelegates(ctx, req.(*GetDelegatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Delegate_ServiceDesc is the grpc.ServiceDesc for Delegate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Delegate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "delegatepb.Delegate",
	HandlerType: (*DelegateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDelegates",
			Handler:    _Delegate_GetDelegates_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "delegatepb/delegate.proto",
}
