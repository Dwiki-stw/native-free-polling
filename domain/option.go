package domain

import (
	"context"
	"native-free-pollings/models"
)

type OptionRepository interface {
	Create(ctx context.Context, db DB, option *models.PollOption) error
	Update(ctx context.Context, db DB, option *models.PollOption) error
	Delete(ctx context.Context, db DB, id int64) error
	GetByPollID(ctx context.Context, db DB, id int64) ([]models.PollOption, error)
}
