package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/goverland-labs/goverland-platform-events/events/ipfs"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
)

const (
	VoteMessage            MessageType = "vote"
	ProposalMessage        MessageType = "proposal"
	InvalidProposalMessage MessageType = "invalid-proposal"
	ArchiveProposalMessage MessageType = "archive-proposal"
	DeleteProposalMessage  MessageType = "delete-proposal"
	SettingsMessage        MessageType = "settings"
)

type MessageType string

type Message struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt  `gorm:"index"`
	MCI       int             `gorm:"unique_index"`
	Space     string          `gorm:"index:space_type_idx"`
	Timestamp time.Time       `gorm:"index:space_type_idx"`
	Type      MessageType     `gorm:"index:space_type_idx"`
	Snapshot  json.RawMessage `gorm:"type:jsonb;serializer:json"`
	IpfsID    string          `gorm:"-"` // virtual property
}

type MessageRepo struct {
	conn *gorm.DB
}

func NewMessageRepo(conn *gorm.DB) *MessageRepo {
	return &MessageRepo{
		conn: conn,
	}
}

func (r *MessageRepo) Upsert(m *Message) (isNew bool, err error) {
	result := r.conn.
		Where(Message{MCI: m.MCI}).
		FirstOrCreate(&m)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *MessageRepo) CreateInBatches(messages []*Message) (newRows int, err error) {
	result := r.conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).CreateInBatches(messages, 500)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (r *MessageRepo) FindLatestMCI() (int, error) {
	var dummy Message
	_ = dummy.MCI

	var mci []int
	err := r.conn.
		Select("mci").
		Table("messages").
		Order("mci desc").
		Limit(1).
		Pluck("mci", &mci).
		Error

	if err != nil {
		return 0, err
	}

	if len(mci) != 1 {
		return 0, nil
	}

	return mci[0], nil
}

func (r *MessageRepo) FindSpacesWithNewVotes(after time.Time) ([]string, error) {
	var dummy Message
	_ = dummy.Space
	_ = dummy.Timestamp

	var spaces []string
	err := r.conn.
		Select("space").
		Distinct().
		Table("messages").
		Where("type = @type", sql.Named("type", VoteMessage)).
		Where("space != ''").
		Where("timestamp >= @timestamp", sql.Named("timestamp", after)).
		Pluck("space", &spaces).
		Error

	if err != nil {
		return nil, err
	}

	return spaces, nil
}

func (r *MessageRepo) FindDeleteProposals(limit, offset int) ([]string, error) {
	var dummy Message
	_ = dummy.Snapshot

	var snapshots []string
	err := r.conn.
		Select("snapshot").
		Table("messages").
		Where("type = @type", sql.Named("type", DeleteProposalMessage)).
		Where("space != ''").
		Where("timestamp >= @timestamp", sql.Named("timestamp", time.Now().Add(-365*24*time.Hour))).
		Limit(limit).
		Offset(offset).
		Pluck("snapshot", &snapshots).
		Error

	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func (r *MessageRepo) FindDaoUpdateSettings(limit, offset int) ([]string, error) {
	var dummy Message
	_ = dummy.Snapshot

	var snapshots []string
	err := r.conn.
		Select("snapshot").
		Table("messages").
		Where("type = @type", sql.Named("type", SettingsMessage)).
		Where("space != ''").
		Where("timestamp >= @timestamp", sql.Named("timestamp", time.Now().Add(-365*24*time.Hour))).
		Limit(limit).
		Offset(offset).
		Pluck("snapshot", &snapshots).
		Error

	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

type MessageService struct {
	repo      *MessageRepo
	publisher Publisher
}

func NewMessageService(repo *MessageRepo, publisher Publisher) *MessageService {
	return &MessageService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *MessageService) Upsert(message ...*Message) error {
	for _, msg := range message {
		msg.Snapshot = helpers.EscapeIllegalCharactersJson(msg.Snapshot)

		// publishing all incoming messages to the ipfs fetcher
		if err := s.publisher.PublishJSON(context.Background(), ipfs.SubjectMessageCreated, ipfs.MessagePayload{
			IpfsID: msg.IpfsID,
			Type:   string(msg.Type),
		}); err != nil {
			log.Error().Err(err).Msg("Failed to publish message")
		}
	}

	_, err := s.repo.CreateInBatches(message)
	if err != nil {
		return err
	}

	return err
}

func (s *MessageService) GetLatestMCI() (int, error) {
	mci, err := s.repo.FindLatestMCI()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	return mci, err
}

func (s *MessageService) FindSpacesWithNewVotes(after time.Time) ([]string, error) {
	spaces, err := s.repo.FindSpacesWithNewVotes(after)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return spaces, err
}

func (s *MessageService) FindDeleteProposals(limit, offset int) ([]string, error) {
	ids, err := s.repo.FindDeleteProposals(limit, offset)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return ids, err
}

func (s *MessageService) FindDaoUpdateSettings(limit, offset int) ([]string, error) {
	ids, err := s.repo.FindDaoUpdateSettings(limit, offset)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return ids, err
}
