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

const spacesPerRequest = 500

type SpacesWorker struct {
	sdk    *snapshot.SDK
	spaces *db.SpaceService

	checkInterval time.Duration
}

func NewSpacesWorker(sdk *snapshot.SDK, spaces *db.SpaceService, checkInterval time.Duration) *SpacesWorker {
	return &SpacesWorker{
		sdk:    sdk,
		spaces: spaces,

		checkInterval: checkInterval,
	}
}

func (w *SpacesWorker) Start(ctx context.Context) error {
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

func (w *SpacesWorker) loop(ctx context.Context) error {
	unknownSpaces, err := w.spaces.GetUndefinedSpaceIDs(spacesPerRequest)
	if err != nil {
		return err
	}
	if len(unknownSpaces) == 0 {
		return nil
	}

	spaces, err := w.fetchSpacesInternal(ctx, []snapshot.ListSpaceOption{snapshot.ListSpaceWithIDs(unknownSpaces...)})
	if err != nil {
		return err
	}

	log.Info().Int("count", len(spaces)).Msg("fetched spaces")

	if err := w.processSpaces(spaces); err != nil {
		return err
	}

	return nil
}

func (w *SpacesWorker) fetchSpacesInternal(ctx context.Context, opts []snapshot.ListSpaceOption) ([]*client.SpaceFragment, error) {
	for {
		spaces, err := w.sdk.ListSpace(ctx, opts...)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}

		return spaces, err
	}
}

func (w *SpacesWorker) processSpaces(spaces []*client.SpaceFragment) error {
	for _, space := range spaces {
		marshaled, err := json.Marshal(space)
		if err != nil {
			return err
		}

		s := db.Space{
			ID:        space.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Snapshot:  marshaled,
		}

		if err := w.spaces.Upsert(&s); err != nil {
			return err
		}

		// FIXME: Send to the queue
	}

	return nil
}
