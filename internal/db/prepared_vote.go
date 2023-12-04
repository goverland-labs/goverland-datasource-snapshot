package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreparedVote struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid"`
	CreatedAt time.Time

	Voter    string
	Proposal string

	TypedData string
}

type PreparedVoteRepo struct {
	conn *gorm.DB
}

func NewPreparedVoteRepo(conn *gorm.DB) *PreparedVoteRepo {
	return &PreparedVoteRepo{
		conn: conn,
	}
}

func (r *PreparedVoteRepo) Create(pv *PreparedVote) error {
	return r.conn.Create(pv).Error
}

func (r *PreparedVoteRepo) Get(id uuid.UUID) (PreparedVote, error) {
	var pv PreparedVote
	err := r.conn.First(&pv, id).Error

	return pv, err
}
