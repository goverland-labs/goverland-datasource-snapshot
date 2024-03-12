package voting

import (
	"context"

	"github.com/google/uuid"
	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/goverland-labs/snapshot-sdk-go/snapshot"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
)

type proposalGetter interface {
	GetByID(id string) (*db.Proposal, error)
}

type snapshotSDK interface {
	Validate(_ context.Context, params snapshot.ValidationParams) (snapshot.ValidationResponse, error)
	GetVotingPower(_ context.Context, params snapshot.GetVotingPowerParams) (*client.VotingPowerFragment, error)
	Vote(_ context.Context, params snapshot.VoteParams) (snapshot.VoteResult, error)
	VoteByID(ctx context.Context, id string) (*client.VoteFragment, error)
}

type preparedVoteStorage interface {
	Create(vote *db.PreparedVote) error
	Get(id uuid.UUID) (db.PreparedVote, error)
}
