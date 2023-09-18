package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/goverland-labs/platform-events/events/aggregator"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/goverland-labs/datasource-snapshot/pkg/communicate"
)

type Vote struct {
	ID           string `gorm:"primarykey"`
	Ipfs         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Voter        string
	SpaceID      string
	ProposalID   string
	Choice       json.RawMessage `gorm:"serializer:json"`
	Reason       string
	App          string
	Vp           float64
	VpByStrategy []float64 `gorm:"serializer:json"`
	VpState      string
	Published    bool
}

type VoteRepo struct {
	conn *gorm.DB
}

func NewVoteRepo(conn *gorm.DB) *VoteRepo {
	return &VoteRepo{
		conn: conn,
	}
}

func (r *VoteRepo) Upsert(v *Vote) (isNew bool, err error) {
	result := r.conn.
		Where(Vote{ID: v.ID}).
		FirstOrCreate(&v)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

// BatchCreate creates votes in batch
func (r *VoteRepo) BatchCreate(data []Vote) error {
	return r.conn.Model(&Vote{}).Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(data, 500).Error
}

func (r *VoteRepo) GetLatestVote(id string) (*Vote, error) {
	var (
		dummy = Vote{}
		_     = dummy.CreatedAt
		_     = dummy.ProposalID
	)

	var v Vote

	err := r.conn.
		Where("proposal_id = ?", id).
		Order("created_at DESC").
		First(&v).
		Error

	return &v, err
}

func (r *VoteRepo) SelectForPublish(limit int) ([]Vote, error) {
	var list []Vote

	err := r.conn.
		Where("published is false").
		Order("created_at ASC").
		Limit(limit).
		Find(&list).
		Error

	return list, err
}

func (r *VoteRepo) MarkAsPublished(votes []Vote) error {
	ids := make([]string, len(votes))
	for i := range votes {
		ids[i] = votes[i].ID
	}

	return r.conn.
		Model(Vote{}).
		Where("id in ?", ids).
		Update("published", true).
		Error
}

type VoteService struct {
	repo      *VoteRepo
	publisher *communicate.Publisher
}

func NewVoteService(repo *VoteRepo, publisher *communicate.Publisher) *VoteService {
	return &VoteService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *VoteService) Upsert(vote *Vote) error {
	_, err := s.repo.Upsert(vote)

	return err
}

func (s *VoteService) BatchCreate(votes []Vote) error {
	return s.repo.BatchCreate(votes)
}

func (s *VoteService) GetLatestVote(id string) (*Vote, error) {
	p, err := s.repo.GetLatestVote(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *VoteService) Publish(limit int) error {
	votes, err := s.repo.SelectForPublish(limit)
	if err != nil {
		return fmt.Errorf("select votes: %w", err)
	}

	log.Info().Msgf("selected %d votes", len(votes))

	if len(votes) == 0 {
		return nil
	}

	pl := make(aggregator.VotesPayload, len(votes))
	for i := range votes {
		pl[i] = aggregator.VotePayload{
			ID:            votes[i].ID,
			Ipfs:          votes[i].Ipfs,
			Voter:         votes[i].Voter,
			Created:       int(votes[i].CreatedAt.Unix()),
			OriginalDaoID: votes[i].SpaceID,
			ProposalID:    votes[i].ProposalID,
			Choice:        votes[i].Choice,
			Reason:        votes[i].Reason,
			App:           votes[i].App,
			Vp:            votes[i].Vp,
			VpByStrategy:  votes[i].VpByStrategy,
			VpState:       votes[i].VpState,
		}
	}

	err = s.publisher.PublishJSON(context.Background(), aggregator.SubjectVoteCreated, pl)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	log.Info().Msg("votes are published")

	err = s.repo.MarkAsPublished(votes)
	if err != nil {
		return fmt.Errorf("mark as published: %w", err)
	}

	log.Info().Msg("votes are marked as read")

	return nil
}
