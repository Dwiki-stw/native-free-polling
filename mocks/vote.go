package mocks

import (
	"context"
	"native-free-pollings/domain"
	"native-free-pollings/models"

	"github.com/stretchr/testify/mock"
)

type VoteRepostoryMock struct {
	mock.Mock
}

func (m *VoteRepostoryMock) Create(ctx context.Context, db domain.DB, vote *models.Vote) error {
	args := m.Called(ctx, db, vote)
	return args.Error(0)
}

func (m *VoteRepostoryMock) CreateUserVote(ctx context.Context, db domain.DB, userID, voteID int64) error {
	args := m.Called(ctx, db, userID, voteID)
	return args.Error(0)
}

func (m *VoteRepostoryMock) GetByPollID(ctx context.Context, db domain.DB, pollID int64) ([]models.Vote, error) {
	args := m.Called(ctx, db, pollID)
	if result, ok := args.Get(0).([]models.Vote); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *VoteRepostoryMock) GetByOptionID(ctx context.Context, db domain.DB, optionID int64) ([]models.Vote, error) {
	args := m.Called(ctx, db, optionID)
	if result, ok := args.Get(0).([]models.Vote); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *VoteRepostoryMock) HasUserVoted(ctx context.Context, db domain.DB, pollID, userID int64) (bool, error) {
	args := m.Called(ctx, db, pollID, userID)
	if result, ok := args.Get(0).(bool); ok {
		return result, args.Error(1)
	}

	return false, args.Error(1)
}

func (m *VoteRepostoryMock) HasDeviceVoted(ctx context.Context, db domain.DB, deviceHash string, pollID int64) (bool, error) {
	args := m.Called(ctx, db, deviceHash, pollID)
	if result, ok := args.Get(0).(bool); ok {
		return result, args.Error(1)
	}

	return false, args.Error(1)
}
