package cli

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/goverland-labs/snapshot-sdk-go/snapshot"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
)

const FetchProposalCommandName = "fetchproposal"

type FetchProposal struct {
	base

	Sdk       *snapshot.SDK
	Spaces    *db.SpaceService
	Proposals *db.ProposalService
}

func (c *FetchProposal) GetName() string {
	return FetchProposalCommandName
}

func (c *FetchProposal) GetArguments() ArgumentsDetails {
	return ArgumentsDetails{
		"space": "proposals space",
	}
}

func (c *FetchProposal) ParseArgs(args ...string) (Arguments, error) {
	return c.parseArgs(c, args...)
}

func (c *FetchProposal) Execute(args Arguments) error {
	start := time.Now()
	log.Info().Msg("FetchProposal started")
	ctx := context.Background()

	space, err := c.getSpace(args)
	if err != nil {
		return err
	}

	err = c.fetchProposals(ctx, space)
	if err != nil {
		return err
	}

	log.Info().Msgf("FetchProposal finished. Took: %v", time.Since(start))

	return nil
}

func (c *FetchProposal) getSpace(args Arguments) (string, error) {
	src := args.Get("space")
	if src == "" {
		return "", errors.New("space is required")
	}

	return src, nil
}

func (c *FetchProposal) fetchProposals(ctx context.Context, space string) error {
	opts := []snapshot.ListProposalOption{snapshot.ListProposalWithOrderBy("created", client.OrderDirectionAsc),
		snapshot.ListProposalWithSpacesFilter(space)}
	proposals, err := c.Sdk.ListProposal(ctx, opts...)
	if err != nil {
		return err
	}
	for _, proposal := range proposals {
		marshaled, err := json.Marshal(proposal)
		if err != nil {
			return err
		}

		p := db.Proposal{
			ID:        proposal.ID,
			SpaceID:   proposal.GetSpace().GetID(),
			CreatedAt: time.Unix(proposal.GetCreated(), 0),
			UpdatedAt: time.Now(),
			Snapshot:  marshaled,
		}

		if err := c.Proposals.Upsert(&p); err != nil {
			return err
		}
	}

	return nil
}
