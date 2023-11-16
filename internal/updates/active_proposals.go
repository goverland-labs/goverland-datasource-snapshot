package updates

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/goverland-labs/sdk-snapshot-go/client"
	"github.com/goverland-labs/sdk-snapshot-go/snapshot"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/datasource-snapshot/internal/db"
)

const gap = 30 * time.Minute

type ActiveProposalsWorker struct {
	sdk       *snapshot.SDK
	proposals *db.ProposalService

	checkInterval time.Duration
}

func NewActiveProposalsWorker(sdk *snapshot.SDK, proposals *db.ProposalService, checkInterval time.Duration) *ActiveProposalsWorker {
	return &ActiveProposalsWorker{
		sdk:       sdk,
		proposals: proposals,

		checkInterval: checkInterval,
	}
}

func (w *ActiveProposalsWorker) Start(ctx context.Context) error {
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

func (w *ActiveProposalsWorker) loop(ctx context.Context) error {
	ids, err := w.proposals.GetProposalIDsForUpdate(gap, proposalsPerRequest, false)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	proposals, err := w.fetchProposalsInternal(ctx, []snapshot.ListProposalOption{
		snapshot.ListProposalWithOrderBy("created", client.OrderDirectionAsc),
		snapshot.ListProposalWithIDFilter(ids...),
	})

	if err != nil {
		return err
	}

	log.Info().Int("requested", len(ids)).Int("fetched", len(proposals)).Msg("updated proposals")

	if err := w.processProposals(proposals); err != nil {
		return err
	}

	forDeletion := findProposalIDsForDeletion(ids, proposals)
	log.Info().Int("delete", len(forDeletion)).Msg("deleted proposals")

	if err := w.proposals.Delete(forDeletion); err != nil {
		return err
	}

	return nil
}

func findProposalIDsForDeletion(requestedIDs []string, fetched []*client.ProposalFragment) []string {
	ids := make([]string, 0, len(requestedIDs))

	for i := range requestedIDs {
		if containsID(requestedIDs[i], fetched) {
			continue
		}

		ids = append(ids, requestedIDs[i])
	}

	return ids
}

func containsID(id string, proposals []*client.ProposalFragment) bool {
	for i := range proposals {
		if proposals[i].ID == id {
			return true
		}
	}

	return false
}

func (w *ActiveProposalsWorker) fetchProposalsInternal(ctx context.Context, opts []snapshot.ListProposalOption) ([]*client.ProposalFragment, error) {
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

func (w *ActiveProposalsWorker) processProposals(proposals []*client.ProposalFragment) error {
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
	}

	return nil
}
