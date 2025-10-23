package mocks

import (
	"context"
	"native-free-pollings/domain"
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
