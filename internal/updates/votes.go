package updates

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/goverland-labs/sdk-snapshot-go/client"
	"github.com/goverland-labs/sdk-snapshot-go/snapshot"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/gorm"

	"github.com/goverland-labs/datasource-snapshot/internal/db"
	"github.com/goverland-labs/datasource-snapshot/internal/helpers"
)

const (
	votesPerRequest   = 1000
	votePublishLimit  = 500
	voteMaxOffset     = 5000
	voteProposalLimit = 20
)

type VoteWorker struct {
	sdk       *snapshot.SDK
	votes     *db.VoteService
	proposals *db.ProposalService

	checkInterval time.Duration
}

func NewVotesWorker(
	sdk *snapshot.SDK,
	votes *db.VoteService,
	proposals *db.ProposalService,
	checkInterval time.Duration,
) *VoteWorker {
	return &VoteWorker{
		sdk:       sdk,
		votes:     votes,
		proposals: proposals,

		checkInterval: checkInterval,
	}
}

func (w *VoteWorker) LoadHistorical(ctx context.Context) error {
	for {
		if err := w.loop(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		case <-time.After(w.checkInterval):
		}
	}
}

func (w *VoteWorker) LoadActive(ctx context.Context) error {
	for {
		if err := w.loopActive(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		case <-time.After(w.checkInterval):
		}
	}
}

func (w *VoteWorker) Publish(ctx context.Context) error {
	for {
		if err := w.votes.Publish(votePublishLimit); err != nil {
			log.Error().Err(err).Msg("vote publish")
		}

		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		case <-time.After(10 * time.Second):
		}
	}
}

func (w *VoteWorker) loop(ctx context.Context) error {
	ids, err := w.proposals.GetProposalForVotes(voteProposalLimit)
	if err != nil {
		return fmt.Errorf("get proposals for votes: %w", err)
	}

	for i := range ids {
		id := ids[i]
		opts := []snapshot.ListVotesOption{
			snapshot.ListVotesWithOrderBy("created", client.OrderDirectionAsc),
			snapshot.ListVotesWithProposalIDsFilter(id),
		}

		createdAfter := w.getLastVoteCreatedAt(id)
		if !createdAfter.IsZero() {
			opts = append(opts, snapshot.ListVotesCreatedAfter(createdAfter))
		}

		offset := 0
		for {
			votes, err := w.fetchVotes(ctx, offset, opts)
			if errors.Is(err, snapshot.ErrTooManyRequests) {
				log.Warn().Err(err).Msg("snapshot api limits are reached")
				<-time.After(5 * time.Second)
				continue
			}
			if err != nil {
				return err
			}

			log.Info().Int("count", len(votes)).Msg("fetched votes")

			if err := w.processVotes(votes); err != nil {
				return err
			}

			if len(votes) < votesPerRequest {
				err := w.proposals.MarkVotesProcessed(id)
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
	}

	return nil
}

func (w *VoteWorker) loopActive(ctx context.Context) error {
	ids, err := w.proposals.GetProposalIDsForUpdate(gap, proposalsPerRequest)
	if err != nil {
		return fmt.Errorf("get proposals for votes: %w", err)
	}

	for i := range ids {
		id := ids[i]
		opts := []snapshot.ListVotesOption{
			snapshot.ListVotesWithOrderBy("created", client.OrderDirectionAsc),
			snapshot.ListVotesWithProposalIDsFilter(id),
		}

		createdAfter := w.getLastVoteCreatedAt(id)
		if !createdAfter.IsZero() {
			opts = append(opts, snapshot.ListVotesCreatedAfter(createdAfter))
		}

		offset := 0
		for {
			votes, err := w.fetchVotes(ctx, offset, opts)
			if errors.Is(err, snapshot.ErrTooManyRequests) {
				log.Warn().Err(err).Msg("snapshot api limits are reached")
				<-time.After(5 * time.Second)
				continue
			}
			if err != nil {
				return err
			}

			log.Info().Int("count", len(votes)).Msg("fetched votes")

			if err := w.processVotes(votes); err != nil {
				return err
			}

			if len(votes) < votesPerRequest {
				break
			}

			offset += votesPerRequest
			if offset > voteMaxOffset {
				break
			}
		}
	}

	return nil
}

func (w *VoteWorker) fetchVotes(ctx context.Context, offset int, opts []snapshot.ListVotesOption) ([]*client.VoteFragment, error) {
	votes, err := w.fetchVotesInternal(ctx, append(opts, snapshot.ListVotesWithPagination(votesPerRequest, offset)))
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
		part, err := w.fetchVotesInternal(ctx, append(opts, snapshot.ListVotesWithPagination(interval.Limit, interval.From)))
		if err != nil {
			return nil, err
		}

		votes = append(votes, part...)
	}

	return votes, nil
}

func (w *VoteWorker) fetchVotesInternal(ctx context.Context, opts []snapshot.ListVotesOption) ([]*client.VoteFragment, error) {
	for {
		votes, err := w.sdk.ListVotes(ctx, opts...)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}

		return votes, err
	}
}

func (w *VoteWorker) processVotes(votes []*client.VoteFragment) error {
	for _, vote := range votes {
		p := db.Vote{
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

		if err := w.votes.Upsert(&p); err != nil {
			return err
		}
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

func (w *VoteWorker) getLastVoteCreatedAt(id string) time.Time {
	var createdAfter time.Time

	lastVote, err := w.votes.GetLatestVote(id)
	if err != nil {
		log.Error().Err(err).Msg("unable to get last fetched proposal")
		return createdAfter
	}

	if lastVote != nil {
		createdAfter = lastVote.CreatedAt.Add(-createdAtGap)
	}

	return createdAfter
}
