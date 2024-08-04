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
	delegates, err := g.service.GetDelegates(ctx, GetDelegatesParams{
		Dao:       req.GetDaoOriginalId(),
		Strategy:  req.GetStrategy().GetValue(),
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

func (g *GrpcServer) GetDelegateProfile(ctx context.Context, req *delegatepb.GetDelegateProfileRequest) (*delegatepb.GetDelegateProfileResponse, error) {
	profile, err := g.service.GetDelegateProfile(ctx, GetDelegateProfileParams{
		Dao:      req.GetDaoOriginalId(),
		Strategy: req.GetStrategy().GetValue(),
		Address:  req.GetAddress(),
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("dao", req.DaoOriginalId).
			Str("address", req.Address).
			Msg("failed to get delegate profile")

		return nil, status.Errorf(codes.Internal, "failed to get delegate profile: %v", err)
	}

	delegates := make([]*delegatepb.ProfileDelegateItem, 0, len(profile.Delegates))
	for _, d := range profile.Delegates {
		delegates = append(delegates, &delegatepb.ProfileDelegateItem{
			Address:         d.Address,
			PercentOfWeight: d.PercentOfWeight,
			DelegatedPower:  d.DelegatedPower,
		})
	}

	return &delegatepb.GetDelegateProfileResponse{
		Address:              profile.Address,
		VotingPower:          profile.VotingPower,
		IncomingPower:        profile.IncomingPower,
		OutgoingPower:        profile.OutgoingPower,
		PercentOfVotingPower: profile.PercentOfVotingPower,
		PercentOfDelegators:  profile.PercentOfDelegators,
		Delegates:            delegates,
	}, nil
}
