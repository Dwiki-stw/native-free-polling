package mocks

import (
	"context"
	"native-free-pollings/domain"
	"native-free-pollings/dto"
	"native-free-pollings/models"

	"github.com/stretchr/testify/mock"
)

type PollRepositoryMock struct {
	mock.Mock
}

func (m *PollRepositoryMock) Create(ctx context.Context, db domain.DB, poll *models.Polling) error {
	args := m.Called(ctx, db, poll)
	return args.Error(0)
}

func (m *PollRepositoryMock) Delete(ctx context.Context, db domain.DB, id int64) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *PollRepositoryMock) Update(ctx context.Context, db domain.DB, poll *models.Polling) error {
	args := m.Called(ctx, db, poll)
	return args.Error(0)
}

func (m *PollRepositoryMock) GetByID(ctx context.Context, db domain.DB, id int64) (*models.Polling, error) {
	args := m.Called(ctx, db, id)
	if poll, ok := args.Get(0).(*models.Polling); ok {
		return poll, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *PollRepositoryMock) GetResultsByID(ctx context.Context, db domain.DB, id int64) ([]models.VoteResult, error) {
	args := m.Called(ctx, db, id)
	if result, ok := args.Get(0).([]models.VoteResult); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

type PollServiceMock struct {
	mock.Mock
}

func (m *PollServiceMock) CreatePolling(ctx context.Context, rq *dto.CreatePollingRequest, creator dto.CreatorInfo) (*dto.PollingResponse, error) {
	args := m.Called(ctx, rq, creator)
	if result, ok := args.Get(0).(*dto.PollingResponse); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *PollServiceMock) UpdatePolling(ctx context.Context, rq *dto.UpdatePollingRequest, creator dto.CreatorInfo) (*dto.PollingResponse, error) {
	args := m.Called(ctx, rq, creator)
	if result, ok := args.Get(0).(*dto.PollingResponse); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *PollServiceMock) VoteOptionPolling(ctx context.Context, userID, pollID, optionID int64, deviceHash string) error {
	args := m.Called(ctx, userID, pollID, optionID, deviceHash)

	return args.Error(0)
}

func (m *PollServiceMock) DeletePolling(ctx context.Context, pollID, userID int64) error {
	args := m.Called(ctx, pollID, userID)

	return args.Error(0)
}

func (m *PollServiceMock) GetDetailPolling(ctx context.Context, id int64) (*dto.PollingResponse, error) {
	args := m.Called(ctx, id)
	if result, ok := args.Get(0).(*dto.PollingResponse); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *PollServiceMock) GetPollingResult(ctx context.Context, pollID int64) (*dto.ResultPolling, error) {
	args := m.Called(ctx, pollID)
	if result, ok := args.Get(0).(*dto.ResultPolling); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}
