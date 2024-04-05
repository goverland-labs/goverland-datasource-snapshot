package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/goverland-labs/goverland-datasource-snapshot/internal/helpers"
)

const (
	VoteMessage            MessageType = "vote"
	ProposalMessage        MessageType = "proposal"
	InvalidProposalMessage MessageType = "invalid-proposal"
	ArchiveProposalMessage MessageType = "archive-proposal"
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

type MessageService struct {
	repo *MessageRepo
}

func NewMessageService(repo *MessageRepo) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (s *MessageService) Upsert(message ...*Message) error {
	for _, msg := range message {
		msg.Snapshot = helpers.EscapeIllegalCharactersJson(msg.Snapshot)
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
