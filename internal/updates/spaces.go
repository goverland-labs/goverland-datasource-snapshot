package updates

import (
	"context"
	"errors"
	"time"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
)

const spacesPerRequest = 100

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
