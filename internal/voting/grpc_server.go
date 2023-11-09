package voting

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/datasource-snapshot/proto/votingpb"
)

type GrpcServer struct {
	votingpb.UnimplementedVotingServer

	actionService *ActionService
}

func NewGrpcServer(actionService *ActionService) *GrpcServer {
	return &GrpcServer{}
}

func (g *GrpcServer) Validate(ctx context.Context, req *votingpb.ValidateRequest) (*votingpb.ValidateResponse, error) {
	if req.GetVoter() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "voter is required")
	}
	if req.GetProposal() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "proposal is required")
	}

	params := &ValidateParams{
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

func (g *GrpcServer) Prepare(ctx context.Context, req *votingpb.PrepareRequest) (*votingpb.PrepareResponse, error) {
	if req.GetVoter() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "voter is required")
	}
	if req.GetProposal() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "proposal is required")
	}

	params := &PrepareParams{
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
		TypedData: prepare.TypedData,
	}, nil
}

func (g *GrpcServer) Vote(ctx context.Context, req *votingpb.VoteRequest) (*votingpb.VoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
