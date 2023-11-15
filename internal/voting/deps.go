package voting

import (
	"context"

	"github.com/goverland-labs/sdk-snapshot-go/snapshot"

	"github.com/goverland-labs/datasource-snapshot/internal/db"
)

type proposalGetter interface {
	GetByID(id string) (*db.Proposal, error)
}

type snapshotSDK interface {
	Validate(_ context.Context, params snapshot.ValidationParams) (snapshot.ValidationResponse, error)
	GetVotingPower(_ context.Context, params snapshot.GetVotingPowerParams) (snapshot.GetVotingPowerResponse, error)
	Vote(_ context.Context, params snapshot.VoteParams) (snapshot.VoteResult, error)
}

type preparedVoteStorage interface {
	Create(vote *db.PreparedVote) error
	Get(id uint64) (db.PreparedVote, error)
}
