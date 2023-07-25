package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/goverland-labs/platform-events/events/aggregator"
	"github.com/goverland-labs/sdk-snapshot-go/client"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/goverland-labs/datasource-snapshot/internal/helpers"
	"github.com/goverland-labs/datasource-snapshot/pkg/communicate"
)

type Proposal struct {
	ID        string `gorm:"primarykey"`
	SpaceID   string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt  `gorm:"index"`
	Snapshot  json.RawMessage `gorm:"type:jsonb;serializer:json"`
}

type ProposalRepo struct {
	conn *gorm.DB
}

func NewProposalRepo(conn *gorm.DB) *ProposalRepo {
	return &ProposalRepo{
		conn: conn,
	}
}

func (r *ProposalRepo) Upsert(p *Proposal) (isNew bool, err error) {
	var existed *Proposal
	err = r.conn.
		Select("id").
		Where(Proposal{ID: p.ID}).
		First(existed).
		Error

	isNew = errors.Is(err, gorm.ErrRecordNotFound)

	p.UpdatedAt = time.Now()
	result := r.conn.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(&p)

	if result.Error != nil {
		return false, result.Error
	}

	return isNew, nil
}

func (r *ProposalRepo) GetLatestProposal() (*Proposal, error) {
	var (
		dummy = Proposal{}
		_     = dummy.CreatedAt
	)

	var p Proposal

	err := r.conn.
		Order("created_at DESC").
		First(&p).
		Error

	return &p, err
}

func (r *ProposalRepo) GetProposalIDsForUpdate(interval time.Duration, limit int) ([]string, error) {
	var (
		dummy = Proposal{}
		_     = dummy.UpdatedAt
		_     = dummy.Snapshot
	)

	var ids []struct {
		ID string
	}

	err := r.conn.Debug().Select("id").
		Table("proposals").
		Where("updated_at < ?", time.Now().Add(-interval)).
		Where(
			r.conn.
				Where("to_timestamp((snapshot->'start')::double precision) <= now() and to_timestamp((snapshot->'end')::double precision) >= now()").
				Or("updated_at < to_timestamp((snapshot->'end')::double precision) and to_timestamp((snapshot->'end')::double precision) < now()"),
		).
		Order("updated_at asc").
		Limit(limit).
		Scan(&ids).
		Error

	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(ids))
	for _, row := range ids {
		result = append(result, row.ID)
	}

	return result, nil
}

type ProposalService struct {
	repo      *ProposalRepo
	publisher *communicate.Publisher
}

func NewProposalService(repo *ProposalRepo, publisher *communicate.Publisher) *ProposalService {
	return &ProposalService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *ProposalService) Upsert(p *Proposal) error {
	// Remove illegal chars from the whole snapshot
	p.Snapshot = helpers.EscapeIllegalCharactersBytes(p.Snapshot)

	isNew, err := s.repo.Upsert(p)
	if err != nil {
		return err
	}

	if !isNew {
		log.Debug().Str("proposal", fmt.Sprintf("%s/%s", p.SpaceID, p.ID)).Msg("proposal updated")
		return s.publishEvent(aggregator.SubjectProposalUpdated, p)
	}

	log.Debug().Str("proposal", fmt.Sprintf("%s/%s", p.SpaceID, p.ID)).Msg("proposal created")
	return s.publishEvent(aggregator.SubjectProposalCreated, p)
}

func (s *ProposalService) publishEvent(subject string, proposal *Proposal) error {
	var unmarshaled client.ProposalFragment
	if err := json.Unmarshal(proposal.Snapshot, &unmarshaled); err != nil {
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

	scores := make([]float32, 0, len(unmarshaled.Scores))
	for i := range unmarshaled.Scores {
		scores = append(scores, float32(helpers.ZeroIfNil(unmarshaled.Scores[i])))
	}

	return s.publisher.PublishJSON(context.Background(), subject, aggregator.ProposalPayload{
		ID:            proposal.ID,
		Ipfs:          helpers.ZeroIfNil(unmarshaled.Ipfs),
		Author:        unmarshaled.Author,
		Created:       int(unmarshaled.Created),
		DaoID:         unmarshaled.GetSpace().GetID(),
		Network:       unmarshaled.Network,
		Symbol:        unmarshaled.Symbol,
		Type:          helpers.ZeroIfNil(unmarshaled.Type),
		Strategies:    strategies,
		Validation:    aggregator.ValidationPayload{},
		Title:         unmarshaled.Title,
		Body:          helpers.ZeroIfNil(unmarshaled.Body),
		Discussion:    unmarshaled.Discussion,
		Choices:       helpers.ResolvePointers(unmarshaled.Choices),
		Start:         int(unmarshaled.Start),
		End:           int(unmarshaled.End),
		Quorum:        unmarshaled.Quorum,
		Privacy:       helpers.ZeroIfNil(unmarshaled.Privacy),
		Snapshot:      helpers.ZeroIfNil(unmarshaled.Snapshot),
		State:         helpers.ZeroIfNil(unmarshaled.State),
		Link:          helpers.ZeroIfNil(unmarshaled.Link),
		App:           helpers.ZeroIfNil(unmarshaled.App),
		Scores:        scores,
		ScoresState:   helpers.ZeroIfNil(unmarshaled.ScoresState),
		ScoresTotal:   float32(helpers.ZeroIfNil(unmarshaled.ScoresTotal)),
		ScoresUpdated: int(helpers.ZeroIfNil(unmarshaled.ScoresUpdated)),
		Votes:         int(helpers.ZeroIfNil(unmarshaled.Votes)),
	})
}

func (s *ProposalService) GetProposalIDsForUpdate(interval time.Duration, limit int) ([]string, error) {
	return s.repo.GetProposalIDsForUpdate(interval, limit)
}

func (s *ProposalService) GetLatestProposal() (*Proposal, error) {
	p, err := s.repo.GetLatestProposal()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}
