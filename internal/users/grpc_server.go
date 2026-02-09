package users

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/goverland-labs/goverland-datasource-snapshot/protocol/userspb"
)

type GrpcServer struct {
	proto.UnimplementedUsersServer

	service *Service
}

func NewGrpcServer(s *Service) *GrpcServer {
	return &GrpcServer{
		service: s,
	}
}

func (g *GrpcServer) GetUsersInfo(ctx context.Context, req *proto.GetUsersInfoRequest) (*proto.GetUsersInfoResponse, error) {
	if len(req.Addresses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "addresses required")
	}

	info, err := g.service.GetUsersInfo(ctx, req.Addresses)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	list := make([]*proto.UserInfo, 0, len(info))
	for _, details := range info {
		list = append(list, &proto.UserInfo{
			Address: details.Address,
			About:   details.About,
		})
	}

	return &proto.GetUsersInfoResponse{
		List: list,
	}, nil
}
