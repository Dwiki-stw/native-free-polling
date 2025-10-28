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
	repo   domain.UserRepository
	hasher helper.PasswordHasher
}

func NewUserService(repo domain.UserRepository, hasher helper.PasswordHasher) domain.UserService {
	return &userService{repo: repo, hasher: hasher}
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

func (u *userService) ChangePassword(ctx context.Context, id int64, password string) error {
	if id <= 0 || password == "" {
		return helper.NewAppError("BAD_REQUEST", "invalid input", nil)
	}

	hashed, err := u.hasher.Hash(password)
	if err != nil {
		return helper.NewAppError("HASH_FAILED", "failed hash password", err)
	}

	err = u.repo.UpdatePassword(ctx, id, hashed)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.NewAppError("NOT_FOUND", "user not found", err)
		}
		return helper.NewAppError("INTERNAL_ERROR", "failed to change password", err)
	}

	return nil
}

func (u *userService) GetUserCreatedPollings(ctx context.Context, id int64) ([]dto.PollingSummaryForCreator, error) {
	if id <= 0 {
		return nil, helper.NewAppError("AUTH_FAILED", "user ID invalid", nil)
	}

	results := []dto.PollingSummaryForCreator{}
	polls, err := u.repo.FindPollingsByID(ctx, id)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed to get pollings", err)
	}

	for _, poll := range polls {
		ps := dto.PollingSummaryForCreator{
			ID:         poll.ID,
			Title:      poll.Title,
			Status:     poll.Status,
			TotalVotes: poll.TotalVotes,
		}
		results = append(results, ps)
	}

	return results, nil
}

func (u *userService) GetUserVotedPollings(ctx context.Context, id int64) ([]dto.PollingSummaryForVoter, error) {
	if id <= 0 {
		return nil, helper.NewAppError("AUTH_FAILED", "user ID invalid", nil)
	}

	results := []dto.PollingSummaryForVoter{}
	polls, err := u.repo.FindPollingsVotedByID(ctx, id)
	if err != nil {
		return nil, helper.NewAppError("DB_ERROR", "failed to get pollings", err)
	}

	for _, poll := range polls {
		ps := dto.PollingSummaryForVoter{
			ID:        poll.ID,
			Title:     poll.Title,
			Status:    poll.Status,
			UserVoted: poll.UserVotedOption,
		}
		results = append(results, ps)
	}

	return results, nil
}
