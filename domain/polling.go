package domain

import (
	"context"
	"native-free-pollings/models"
)

type PollRepository interface {
	Create(ctx context.Context, db DB, poll *models.Polling) error
	Update(ctx context.Context, db DB, poll *models.Polling) error
	Delete(ctx context.Context, db DB, id int64) error
	GetByID(ctx context.Context, db DB, id int64) (*models.Polling, error)
	GetResultsByID(ctx context.Context, db DB, id int64) ([]models.VoteResult, error)
}
