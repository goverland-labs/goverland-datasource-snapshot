package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/goverland-labs/snapshot-sdk-go/snapshot"
	"github.com/samber/lo"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/gorm"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
)

const (
	FetchCommandName = "fetch"
	votesPerRequest  = 1000
	votePublishLimit = 350
	voteMaxOffset    = 5000
)

type FetchForSpace struct {
	base

	Sdk       *snapshot.SDK
	Spaces    *db.SpaceService
	Proposals *db.ProposalService
	Votes     *db.VoteService
}

func (c *FetchForSpace) GetName() string {
	return FetchCommandName
}

func (c *FetchForSpace) GetArguments() ArgumentsDetails {
	return ArgumentsDetails{
		"space": "space to fetch",
	}
}

func (c *FetchForSpace) ParseArgs(args ...string) (Arguments, error) {
	return c.parseArgs(c, args...)
}

func (c *FetchForSpace) Execute(args Arguments) error {
	start := time.Now()
	log.Info().Msg("FetchForSpace started")
	ctx := context.Background()

	space, err := c.getSpace(args)
	if err != nil {
		return err
	}

	err = c.fetchSpace(ctx, space)
	if err != nil {
		return err
	}

	proposalIds, err := c.fetchProposals(ctx, space)
	if err != nil {
		return err
	}

	for i := range proposalIds {
		if err := c.fetchProposalVotes(ctx, proposalIds[i]); err != nil {
			return err
		}
	}

	log.Info().Msgf("FetchForSpace finished. Took: %v", time.Since(start))

	return nil
}

func (c *FetchForSpace) getSpace(args Arguments) (string, error) {
	src := args.Get("space")
	if src == "" {
		return "", errors.New("space is required")
	}

	return src, nil
}

func (c *FetchForSpace) fetchSpace(ctx context.Context, space string) error {
	s, err := c.Sdk.GetSpaceByID(ctx, space)
	if err != nil {
		return err
	}
	marshaled, err := json.Marshal(s)
	if err != nil {
		return err
	}

	sd := db.Space{
		ID:        s.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Snapshot:  marshaled,
	}

	if err := c.Spaces.Upsert(&sd); err != nil {
		return err
	}

	return nil
}

func (c *FetchForSpace) fetchProposals(ctx context.Context, space string) ([]string, error) {
	opts := []snapshot.ListProposalOption{snapshot.ListProposalWithOrderBy("created", client.OrderDirectionAsc),
		snapshot.ListProposalWithSpacesFilter(space)}
	proposals, err := c.Sdk.ListProposal(ctx, opts...)
	if err != nil {
		return nil, err
	}
	proposalIds := make([]string, 0)
	for _, proposal := range proposals {
		marshaled, err := json.Marshal(proposal)
		if err != nil {
			return proposalIds, err
		}

		p := db.Proposal{
			ID:        proposal.ID,
			SpaceID:   proposal.GetSpace().GetID(),
			CreatedAt: time.Unix(proposal.GetCreated(), 0),
			UpdatedAt: time.Now(),
			Snapshot:  marshaled,
		}

		if err := c.Proposals.Upsert(&p); err != nil {
			return proposalIds, err
		}
		proposalIds = append(proposalIds, proposal.ID)
	}

	return proposalIds, nil
}

func (c *FetchForSpace) fetchProposalVotes(ctx context.Context, id string) error {

	opts := []snapshot.ListVotesOption{
		snapshot.ListVotesWithOrderBy("created", client.OrderDirectionAsc),
		snapshot.ListVotesWithProposalIDsFilter(id),
	}

	offset := 0
	for {
		votes, err := c.fetchVotes(ctx, offset, opts)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}
		if err != nil {
			return err
		}
		if err := c.processVotes(votes); err != nil {
			return err
		}

		if len(votes) < votesPerRequest {
			err := c.Proposals.MarkVotesProcessed(id)
			if err != nil {
				log.Warn().Err(err).Msgf("mark votes processed: %s", id)
			}

			break
		}

		offset += votesPerRequest
		if offset > voteMaxOffset {
			break
		}
	}

	return nil
}

func (c *FetchForSpace) fetchVotes(ctx context.Context, offset int, opts []snapshot.ListVotesOption) ([]*client.VoteFragment, error) {
	votes, err := c.fetchVotesInternal(ctx, append(opts, snapshot.ListVotesWithPagination(votesPerRequest, offset)))
	if err == nil {
		return votes, nil
	}

	gqlErr, ok := err.(*clientv2.ErrorResponse)
	if !ok {
		return votes, err
	}

	if gqlErr.GqlErrors == nil {
		return votes, err
	}

	skipOffsets := make([]int, 0, len(*gqlErr.GqlErrors))
	for _, e := range *gqlErr.GqlErrors {
		if len(e.Path) < 2 {
			continue
		}

		index, ok := e.Path[1].(ast.PathIndex)
		if !ok {
			continue
		}

		skipOffsets = append(skipOffsets, int(index+1))
	}

	intervals := helpers.GenerateIntervals(offset, votesPerRequest, lo.Uniq(skipOffsets))

	votes = make([]*client.VoteFragment, 0, votesPerRequest)
	for _, interval := range intervals {
		part, err := c.fetchVotesInternal(ctx, append(opts, snapshot.ListVotesWithPagination(interval.Limit, interval.From)))
		if err != nil {
			return nil, err
		}

		votes = append(votes, part...)
	}

	return votes, nil
}

func (c *FetchForSpace) fetchVotesInternal(ctx context.Context, opts []snapshot.ListVotesOption) ([]*client.VoteFragment, error) {
	for {
		votes, err := c.Sdk.ListVotes(ctx, opts...)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}

		return votes, err
	}
}

func (c *FetchForSpace) processVotes(votes []*client.VoteFragment) error {
	converted := make([]db.Vote, len(votes))
	for i, vote := range votes {
		converted[i] = db.Vote{
			ID:           vote.ID,
			Ipfs:         *vote.GetIpfs(),
			CreatedAt:    time.Unix(vote.GetCreated(), 0),
			UpdatedAt:    time.Now(),
			DeletedAt:    gorm.DeletedAt{},
			Voter:        vote.GetVoter(),
			SpaceID:      vote.GetSpace().GetID(),
			ProposalID:   vote.GetProposal().GetID(),
			Choice:       vote.GetChoice(),
			Reason:       *vote.GetReason(),
			App:          *vote.GetApp(),
			Vp:           *vote.GetVp(),
			VpByStrategy: convertVpByStrategy(vote.GetVpByStrategy()),
			VpState:      *vote.GetVpState(),
		}
	}

	err := c.Votes.BatchCreate(converted)
	if err != nil {
		return fmt.Errorf("batchCreate: %w", err)
	}

	err = c.Votes.Publish(converted, votePublishLimit)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

func convertVpByStrategy(data []*float64) []float64 {
	res := make([]float64, len(data))
	for i := range data {
		res[i] = *data[i]
	}

	return res
}
