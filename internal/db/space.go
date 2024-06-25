package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/goverland-labs/goverland-platform-events/events/aggregator"
	"github.com/goverland-labs/snapshot-sdk-go/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
)

type Space struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt  `gorm:"index"`
	Snapshot  json.RawMessage `gorm:"type:jsonb;serializer:json"`
}

type SpaceRepo struct {
	conn *gorm.DB
}

func NewSpaceRepo(conn *gorm.DB) *SpaceRepo {
	return &SpaceRepo{
		conn: conn,
	}
}

func (r *SpaceRepo) Upsert(s *Space) (isNew bool, err error) {
	var existed Space

	err = r.conn.
		Select("id").
		Where(Space{ID: s.ID}).
		First(&existed).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	isNew = errors.Is(err, gorm.ErrRecordNotFound)

	s.UpdatedAt = time.Now()
	result := r.conn.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(&s)

	if result.Error != nil {
		return false, result.Error
	}

	return isNew, nil
}

func (r *SpaceRepo) FindUndefinedSpaceIDs(limit int) ([]string, error) {
	var (
		proposal = Proposal{}
		space    = Space{}
		_        = proposal.SpaceID
		_        = space.ID
	)

	var ids []string
	err := r.conn.Select("space_id").
		Distinct().
		Table("proposals").
		Where("space_id NOT IN (?)", r.conn.Select("id").Table("spaces").Where("id = proposals.space_id")).
		Limit(limit).
		Scan(&ids).
		Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}

type SpaceService struct {
	repo      *SpaceRepo
	publisher Publisher
}

func NewSpaceService(repo *SpaceRepo, publisher Publisher) *SpaceService {
	return &SpaceService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *SpaceService) Upsert(space *Space) error {
	space.Snapshot = helpers.EscapeIllegalCharactersJson(space.Snapshot)

	isNew, err := s.repo.Upsert(space)
	if err != nil {
		return err
	}

	if !isNew {
		log.Debug().Str("space", space.ID).Msg("space updated")
		return s.publishEvent(aggregator.SubjectDaoUpdated, space)
	}

	log.Debug().Str("space", space.ID).Msg("space created")
	return s.publishEvent(aggregator.SubjectDaoCreated, space)
}

func (s *SpaceService) publishEvent(subject string, space *Space) error {
	var unmarshaled client.SpaceFragment
	if err := json.Unmarshal(space.Snapshot, &unmarshaled); err != nil {
		return err
	}

	strategies := make([]aggregator.StrategyPayload, 0, len(unmarshaled.Strategies))
	for _, strategy := range unmarshaled.GetStrategies() {
		strategies = append(strategies, aggregator.StrategyPayload{
			Name:    strategy.GetName(),
			Network: helpers.ZeroIfNil(strategy.GetNetwork()),
			Params:  strategy.GetParams(),
		})
	}

	treasuries := make([]aggregator.TreasuryPayload, 0, len(unmarshaled.Strategies))
	for _, treasury := range unmarshaled.GetTreasuries() {
		treasuries = append(treasuries, aggregator.TreasuryPayload{
			Name:    helpers.ZeroIfNil(treasury.GetName()),
			Address: helpers.ZeroIfNil(treasury.GetAddress()),
			Network: helpers.ZeroIfNil(treasury.GetNetwork()),
		})
	}

	return s.publisher.PublishJSON(context.Background(), subject, aggregator.DaoPayload{
		ID:             space.ID,
		Name:           helpers.ZeroIfNil(unmarshaled.Name),
		About:          helpers.ZeroIfNil(unmarshaled.About),
		Avatar:         helpers.ZeroIfNil(unmarshaled.Avatar),
		Terms:          helpers.ZeroIfNil(unmarshaled.Terms),
		Location:       helpers.ZeroIfNil(unmarshaled.Location),
		Website:        helpers.ZeroIfNil(unmarshaled.Website),
		Twitter:        helpers.ZeroIfNil(unmarshaled.Twitter),
		Github:         helpers.ZeroIfNil(unmarshaled.Github),
		Coingecko:      helpers.ZeroIfNil(unmarshaled.Coingecko),
		Email:          helpers.ZeroIfNil(unmarshaled.Email),
		Network:        helpers.ZeroIfNil(unmarshaled.Network),
		Symbol:         helpers.ZeroIfNil(unmarshaled.Symbol),
		Skin:           helpers.ZeroIfNil(unmarshaled.Skin),
		Domain:         helpers.ZeroIfNil(unmarshaled.Domain),
		Admins:         helpers.ResolvePointers(unmarshaled.Admins),
		Members:        helpers.ResolvePointers(unmarshaled.Members),
		Moderators:     helpers.ResolvePointers(unmarshaled.Moderators),
		Voting:         aggregator.VotingPayload{},
		Categories:     helpers.ResolvePointers(unmarshaled.Categories),
		Validation:     aggregator.ValidationPayload{},
		VoteValidation: aggregator.ValidationPayload{},
		FollowersCount: int(helpers.ZeroIfNil(unmarshaled.FollowersCount)),
		ProposalsCount: int(helpers.ZeroIfNil(unmarshaled.ProposalsCount)),
		Guidelines:     helpers.ZeroIfNil(unmarshaled.Guidelines),
		Template:       helpers.ZeroIfNil(unmarshaled.Template),
		Strategies:     strategies,
		Treasures:      treasuries,
		Verified:       helpers.ZeroIfNil(unmarshaled.Verified),
	})
}

func (s *SpaceService) GetUndefinedSpaceIDs(limit int) ([]string, error) {
	return s.repo.FindUndefinedSpaceIDs(limit)
}
