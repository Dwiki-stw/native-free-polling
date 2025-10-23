package domain

import (
	"context"
	"native-free-pollings/dto"
	"native-free-pollings/models"
)

type PollRepository interface {
	Create(ctx context.Context, db DB, poll *models.Polling) error
	Update(ctx context.Context, db DB, poll *models.Polling) error
	Delete(ctx context.Context, db DB, id int64) error
	GetByID(ctx context.Context, db DB, id int64) (*models.Polling, error)
	GetResultsByID(ctx context.Context, db DB, id int64) ([]models.VoteResult, error)
}

type PollService interface {
	CreatePolling(ctx context.Context, rq *dto.CreatePollingRequest, creator dto.CreatorInfo) (*dto.PollingResponse, error)
	UpdatePolling(ctx context.Context, rq *dto.UpdatePollingRequest, creator dto.CreatorInfo) (*dto.PollingResponse, error)
	DeletePolling(ctx context.Context, pollID, userID int64) error
	GetDetailPolling(ctx context.Context, id int64) (*dto.PollingResponse, error)
	VoteOptionPolling(ctx context.Context, userID, pollID, optionID int64, deviceHash string) error
	GetPollingResult(ctx context.Context, pollID int64) (*dto.ResultPolling, error)
}
