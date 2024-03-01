package updates

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/goverland-labs/snapshot-sdk-go/snapshot"

	"github.com/goverland-labs/datasource-snapshot/internal/db"
	"github.com/goverland-labs/datasource-snapshot/internal/helpers"
)

const (
	proposalCreatedAtGap = time.Hour
	proposalsPerRequest  = 500
	proposalsMaxOffset   = 5000
)

type ProposalWorker struct {
	sdk       *snapshot.SDK
	proposals *db.ProposalService

	checkInterval time.Duration
	createdAfter  time.Time
}

func NewProposalsWorker(sdk *snapshot.SDK, proposals *db.ProposalService, checkInterval time.Duration) *ProposalWorker {
	return &ProposalWorker{
		sdk:       sdk,
		proposals: proposals,

		checkInterval: checkInterval,
	}
}

func (w *ProposalWorker) Start(ctx context.Context) error {
	w.createdAfter = w.getLastProposalCreatedAt()

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

func (w *ProposalWorker) loop(ctx context.Context) error {
	opts := []snapshot.ListProposalOption{snapshot.ListProposalWithOrderBy("created", client.OrderDirectionAsc)}
	if !w.createdAfter.IsZero() {
		opts = append(opts, snapshot.ListProposalCreatedAfter(w.createdAfter))
	}

	offset := 0
	for {
		proposals, err := w.fetchProposals(ctx, offset, opts)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}
		if err != nil {
			return err
		}

		log.Info().Int("count", len(proposals)).Msg("fetched proposals")

		if err := w.processProposals(proposals); err != nil {
			return err
		}

		if len(proposals) < proposalsPerRequest {
			break
		}

		offset += proposalsPerRequest
		if offset > proposalsMaxOffset {
			break
		}
	}

	return nil
}

func (w *ProposalWorker) fetchProposals(ctx context.Context, offset int, opts []snapshot.ListProposalOption) ([]*client.ProposalFragment, error) {
	proposals, err := w.fetchProposalsInternal(ctx, append(opts, snapshot.ListProposalWithPagination(proposalsPerRequest, offset)))
	if err == nil {
		return proposals, nil
	}

	gqlErr, ok := err.(*clientv2.ErrorResponse)
	if !ok {
		return proposals, err
	}

	if gqlErr.GqlErrors == nil {
		return proposals, err
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

	intervals := helpers.GenerateIntervals(offset, proposalsPerRequest, lo.Uniq(skipOffsets))

	proposals = make([]*client.ProposalFragment, 0, proposalsPerRequest)
	for _, interval := range intervals {
		part, err := w.fetchProposalsInternal(ctx, append(opts, snapshot.ListProposalWithPagination(interval.Limit, interval.From)))
		if err != nil {
			return nil, err
		}

		proposals = append(proposals, part...)
	}

	return proposals, nil
}

func (w *ProposalWorker) fetchProposalsInternal(ctx context.Context, opts []snapshot.ListProposalOption) ([]*client.ProposalFragment, error) {
	for {
		proposals, err := w.sdk.ListProposal(ctx, opts...)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}

		return proposals, err
	}
}

func (w *ProposalWorker) processProposals(proposals []*client.ProposalFragment) error {
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

		if err := w.proposals.Upsert(&p); err != nil {
			return err
		}

		createdAfter := p.CreatedAt.Add(-proposalCreatedAtGap)
		if w.createdAfter.Before(createdAfter) {
			w.createdAfter = createdAfter
		}
	}

	return nil
}

func (w *ProposalWorker) getLastProposalCreatedAt() time.Time {
	var createdAfter time.Time

	lastProposal, err := w.proposals.GetLatestProposal()
	if err != nil {
		log.Error().Err(err).Msg("unable to get last fetched proposal")
		return createdAfter
	}

	if lastProposal != nil {
		createdAfter = lastProposal.CreatedAt.Add(-proposalCreatedAtGap)
	}

	return createdAfter
}
