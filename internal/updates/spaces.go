package updates

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/goverland-labs/snapshot-sdk-go/client"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
)

const spacesPerRequest = 500

type SpacesWorker struct {
	spaceUpdater *SpacesUpdater
	spaces       *db.SpaceService

	checkInterval time.Duration
}

func NewSpacesWorker(spaceUpdater *SpacesUpdater, spaces *db.SpaceService, checkInterval time.Duration) *SpacesWorker {
	return &SpacesWorker{
		spaceUpdater:  spaceUpdater,
		spaces:        spaces,
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

	return w.spaceUpdater.ProcessSpaces(ctx, unknownSpaces)
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
	}

	return nil
}
