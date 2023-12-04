package grpcsrv

import (
	"context"

	"google.golang.org/grpc"
)

const reflectionMethod = "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo"

func UnaryReflectionFilter(in grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == reflectionMethod {
			return handler(ctx, req)
		}
		return in(ctx, req, info, handler)
	}
}

func StreamReflectionFilter(in grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info.FullMethod == reflectionMethod {
			return handler(srv, ss)
		}
		return in(srv, ss, info, handler)
	}
}
