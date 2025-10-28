package domain

import (
	"context"
	"native-free-pollings/dto"
	"native-free-pollings/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, id int64, passwordHashed string) error
	FindPollingsByID(ctx context.Context, id int64) ([]models.PollingSummary, error)
	FindPollingsVotedByID(ctx context.Context, id int64) ([]models.PollingSummary, error)
}

type UserService interface {
	GetProfile(ctx context.Context, id int64) (*dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, user *models.User) (*dto.ProfileResponse, error)
	ChangePassword(ctx context.Context, id int64, password string) error
}
