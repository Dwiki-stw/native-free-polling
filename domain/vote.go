package domain

import (
	"context"
	"native-free-pollings/models"
)

type VoteRepository interface {
	Create(ctx context.Context, db DB, vote *models.Vote) error
	CreateUserVote(ctx context.Context, db DB, userID, voteID int64) error
	GetByPollID(ctx context.Context, db DB, pollID int64) ([]models.Vote, error)
	GetByOptionID(ctx context.Context, db DB, optionID int64) ([]models.Vote, error)
	HasUserVoted(ctx context.Context, db DB, pollID, userID int64) (bool, error)
	HasDeviceVoted(ctx context.Context, db DB, deviceHash string, pollID int64) (bool, error)
}
