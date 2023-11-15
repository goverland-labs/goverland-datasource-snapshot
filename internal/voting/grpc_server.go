package voting

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/datasource-snapshot/proto/votingpb"
)

type GrpcServer struct {
	votingpb.UnimplementedVotingServer

	actionService *ActionService
}

func NewGrpcServer(actionService *ActionService) *GrpcServer {
	return &GrpcServer{
		actionService: actionService,
	}
}

func (g *GrpcServer) Validate(_ context.Context, req *votingpb.ValidateRequest) (*votingpb.ValidateResponse, error) {
	if strings.TrimSpace(req.GetVoter()) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "voter is required")
	}
	if strings.TrimSpace(req.GetProposal()) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "proposal is required")
	}

	params := ValidateParams{
		Voter:    req.GetVoter(),
		Proposal: req.GetProposal(),
	}

	validate, err := g.actionService.Validate(params)
	if err != nil {
		return nil, err
	}

	var validationError *votingpb.ValidationError
	if validate.ValidationError != nil {
		validationError = &votingpb.ValidationError{
			Message: validate.ValidationError.Message,
			Code:    validate.ValidationError.Code,
		}
	}

	return &votingpb.ValidateResponse{
		Ok:              validate.OK,
		VotingPower:     validate.VotingPower,
		ValidationError: validationError,
	}, nil
}

func (g *GrpcServer) Prepare(_ context.Context, req *votingpb.PrepareRequest) (*votingpb.PrepareResponse, error) {
	if strings.TrimSpace(req.GetVoter()) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "voter is required")
	}
	if strings.TrimSpace(req.GetProposal()) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "proposal is required")
	}

	params := PrepareParams{
		Voter:    req.GetVoter(),
		Proposal: req.GetProposal(),
		Choice:   req.GetChoice().GetValue(),
		Reason:   req.Reason,
	}

	prepare, err := g.actionService.Prepare(params)
	if err != nil {
		return nil, err
	}

	return &votingpb.PrepareResponse{
		Id:        prepare.ID,
		TypedData: prepare.TypedData,
	}, nil
}

func (g *GrpcServer) Vote(_ context.Context, req *votingpb.VoteRequest) (*votingpb.VoteResponse, error) {
	if strings.TrimSpace(req.GetSig()) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "sig is required")
	}

	params := VoteParams{
		ID:  req.GetId(),
		Sig: req.GetSig(),
	}

	voteResult, err := g.actionService.Vote(params)
	if err != nil {
		return nil, err
	}

	return &votingpb.VoteResponse{
		Id:   voteResult.ID,
		Ipfs: voteResult.IPFS,
		Relayer: &votingpb.Relayer{
			Address: voteResult.Relayer.Address,
			Receipt: voteResult.Relayer.Receipt,
		},
	}, nil
}
