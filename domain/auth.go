package domain

import (
	"context"
	"native-free-pollings/models"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}
