package domain

import (
	"context"
	"native-free-pollings/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, id int64, passwordHashed string) error
}
