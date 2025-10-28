package mocks

import (
	"context"
	"native-free-pollings/dto"
	"native-free-pollings/models"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}

	return user.(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdatePassword(ctx context.Context, id int64, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindPollingsByID(ctx context.Context, id int64) ([]models.PollingSummary, error) {
	args := m.Called(ctx, id)
	if results, ok := args.Get(0).([]models.PollingSummary); ok {
		return results, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *UserRepositoryMock) FindPollingsVotedByID(ctx context.Context, id int64) ([]models.PollingSummary, error) {
	args := m.Called(ctx, id)
	if results, ok := args.Get(0).([]models.PollingSummary); ok {
		return results, args.Error(1)
	}

	return nil, args.Error(1)
}

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) GetProfile(ctx context.Context, id int64) (*dto.ProfileResponse, error) {
	args := m.Called(ctx, id)
	if resp, ok := args.Get(0).(*dto.ProfileResponse); ok {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *UserServiceMock) UpdateProfile(ctx context.Context, user *models.User) (*dto.ProfileResponse, error) {
	args := m.Called(ctx, user)
	if resp, ok := args.Get(0).(*dto.ProfileResponse); ok {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *UserServiceMock) ChangePassword(ctx context.Context, id int64, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *UserServiceMock) GetUserCreatedPollings(ctx context.Context, id int64) ([]dto.PollingSummaryForCreator, error) {
	args := m.Called(ctx, id)
	if results, ok := args.Get(0).([]dto.PollingSummaryForCreator); ok {
		return results, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *UserServiceMock) GetUserVotedPollings(ctx context.Context, id int64) ([]dto.PollingSummaryForVoter, error) {
	args := m.Called(ctx, id)
	if results, ok := args.Get(0).([]dto.PollingSummaryForVoter); ok {
		return results, args.Error(1)
	}

	return nil, args.Error(1)
}
