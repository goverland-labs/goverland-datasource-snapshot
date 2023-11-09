package voting

import (
	"context"

	"github.com/goverland-labs/datasource-snapshot/internal/db"
	"github.com/goverland-labs/sdk-snapshot-go/snapshot"
)

type proposalGetter interface {
	GetByID(id string) (*db.Proposal, error)
}

type snapshotSDK interface {
	Validate(_ context.Context, params snapshot.ValidationParams) (snapshot.ValidationResponse, error)
	GetVotingPower(_ context.Context, params snapshot.GetVotingPowerParams) (snapshot.GetVotingPowerResponse, error)
}
