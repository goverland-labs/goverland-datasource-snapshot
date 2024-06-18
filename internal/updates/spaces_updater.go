package updates

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/goverland-labs/snapshot-sdk-go/snapshot"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/db"
)

type SpacesUpdater struct {
	sdk    *snapshot.SDK
	spaces *db.SpaceService
}

func NewSpacesUpdater(sdk *snapshot.SDK, spaces *db.SpaceService) *SpacesUpdater {
	return &SpacesUpdater{
		sdk:    sdk,
		spaces: spaces,
	}
}

func (s *SpacesUpdater) ProcessSpaces(ctx context.Context, spaces []string) error {
	spacesResp, err := s.getSpaces(ctx, []snapshot.ListSpaceOption{snapshot.ListSpaceWithIDFilter(spaces...)})
	if err != nil {
		return fmt.Errorf("list spaces: %w", err)
	}

	for _, space := range spacesResp {
		if err := s.saveSpace(space); err != nil {
			return fmt.Errorf("process space: %w", err)
		}
	}

	return nil
}

func (s *SpacesUpdater) ProcessSpace(ctx context.Context, space string) error {
	spaceResp, err := s.sdk.GetSpaceByID(ctx, space)
	if err != nil {
		return fmt.Errorf("get space: %w", err)
	}

	if err := s.saveSpace(spaceResp); err != nil {
		return fmt.Errorf("process space: %w", err)
	}

	return nil
}

// TODO: rework handling err too many requests
func (s *SpacesUpdater) getSpaces(ctx context.Context, opts []snapshot.ListSpaceOption) ([]*client.SpaceFragment, error) {
	for {
		spaces, err := s.sdk.ListSpace(ctx, opts...)
		if errors.Is(err, snapshot.ErrTooManyRequests) {
			log.Warn().Err(err).Msg("snapshot api limits are reached")
			<-time.After(5 * time.Second)
			continue
		}

		return spaces, err
	}
}

func (s *SpacesUpdater) saveSpace(space *client.SpaceFragment) error {
	marshaled, err := json.Marshal(space)
	if err != nil {
		return err
	}

	spaceEntity := db.Space{
		ID:        space.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Snapshot:  marshaled,
	}

	if err := s.spaces.Upsert(&spaceEntity); err != nil {
		return err
	}

	return nil
}
