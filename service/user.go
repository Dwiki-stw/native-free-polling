package service

import (
	"context"
	"database/sql"
	"errors"
	"native-free-pollings/domain"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"native-free-pollings/models"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (u *userService) GetProfile(ctx context.Context, id int64) (*dto.ProfileResponse, error) {
	resp, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helper.NewAppError("NOT_FOUND", "user not found", err)
		}
		return nil, helper.NewAppError("INTERNAL_ERROR", "failed to get user", err)
	}

	return &dto.ProfileResponse{
		ID:        resp.ID,
		Name:      resp.Name,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}

func (u *userService) UpdateProfile(ctx context.Context, user *models.User) (*dto.ProfileResponse, error) {
	if err := u.repo.Update(ctx, user); err != nil {
		if err == sql.ErrNoRows {
			return nil, helper.NewAppError("NOT_FOUND", "user not found", err)
		}
		return nil, helper.NewAppError("INTERNAL_ERROR", "failed to update profile", err)
	}

	user, err := u.repo.GetByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewAppError("NOT_FOUND", "user not found", err)
		}
		return nil, helper.NewAppError("INTERNAL_ERROR", "failed to get user", err)
	}

	return &dto.ProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil

}

func (u *userService) ChangePassword(ctx context.Context, id int64, passwordHashed string) error {
	if id <= 0 || passwordHashed == "" {
		return helper.NewAppError("BAD_REQUEST", "invalid input", nil)
	}

	err := u.repo.UpdatePassword(ctx, id, passwordHashed)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.NewAppError("NOT_FOUND", "user not found", err)
		}
		return helper.NewAppError("INTERNAL_ERROR", "failed to change password", err)
	}

	return nil
}
