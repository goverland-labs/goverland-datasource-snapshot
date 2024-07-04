package delegate

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"
)

type GrpcServer struct {
	delegatepb.UnimplementedDelegateServer

	service *Service
}

func NewGrpcServer(s *Service) *GrpcServer {
	return &GrpcServer{
		service: s,
	}
}

func (g *GrpcServer) GetDelegates(ctx context.Context, req *delegatepb.GetDelegatesRequest) (*delegatepb.GetDelegatesResponse, error) {
	strategy := req.GetStrategy().GetValue()
	delegates, err := g.service.GetDelegates(GetDelegatesRequest{
		Dao:       req.GetDaoOriginalId(),
		Strategy:  strategy,
		By:        req.GetSort(),
		Addresses: req.GetAddresses(),
		Limit:     int(req.GetLimit()),
		Offset:    int(req.GetOffset()),
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("dao", req.DaoOriginalId).
			Msg("failed to get delegates")

		return nil, status.Errorf(codes.Internal, "failed to get delegates: %v", err)
	}

	delegatesResult := make([]*delegatepb.DelegateEntry, 0, len(delegates))
	for _, d := range delegates {
		delegatesResult = append(delegatesResult, &delegatepb.DelegateEntry{
			Address:                  d.Address,
			DelegatorCount:           d.DelegatorCount,
			PercentOfDelegators:      d.PercentOfDelegators,
			VotingPower:              d.VotingPower,
			PercentOfVotingPower:     d.PercentOfVotingPower,
			About:                    "test about",
			Statement:                "test statement",
			UserDelegatedVotingPower: 0,
		})
	}

	return &delegatepb.GetDelegatesResponse{
		Delegates: delegatesResult,
	}, nil
}
