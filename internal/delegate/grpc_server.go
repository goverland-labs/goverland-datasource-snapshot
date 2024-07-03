package delegate

import (
	"context"

	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"
)

type GrpcServer struct {
	delegatepb.UnimplementedDelegateServer
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (g *GrpcServer) GetDelegates(ctx context.Context, req *delegatepb.GetDelegatesRequest) (*delegatepb.GetDelegatesResponse, error) {
	return &delegatepb.GetDelegatesResponse{
		Delegates: []*delegatepb.DelegateEntry{
			{
				Address:                  "0x952D069AEf7cd1358B44da3118154240aFF99aFF",
				DelegatorCount:           1,
				PercentOfDelegators:      100,
				VotingPower:              123.231,
				PercentOfVotingPower:     12,
				About:                    "about",
				Statement:                "statement",
				UserDelegatedVotingPower: 0,
			},
			{
				Address:                  "0x7697cAB0e123c68d27d7D5A9EbA346d7584Af888",
				DelegatorCount:           2,
				PercentOfDelegators:      300,
				VotingPower:              1353.343,
				PercentOfVotingPower:     1222,
				About:                    "about 2",
				Statement:                "statement 2",
				UserDelegatedVotingPower: 102.2,
			},
		},
	}, nil
}
