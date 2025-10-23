package mocks

import (
	"context"
	"native-free-pollings/domain"
	"native-free-pollings/models"

	"github.com/stretchr/testify/mock"
)

type OptionRepositoryMock struct {
	mock.Mock
}

func (m *OptionRepositoryMock) Create(ctx context.Context, db domain.DB, option *models.PollOption) error {
	args := m.Called(ctx, db, option)
	return args.Error(0)
}

func (m *OptionRepositoryMock) Update(ctx context.Context, db domain.DB, option *models.PollOption) error {
	args := m.Called(ctx, db, option)
	return args.Error(0)
}

func (m *OptionRepositoryMock) Delete(ctx context.Context, db domain.DB, id int64) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *OptionRepositoryMock) GetByPollID(ctx context.Context, db domain.DB, id int64) ([]models.PollOption, error) {
	args := m.Called(ctx, db, id)
	if result, ok := args.Get(0).([]models.PollOption); ok {
		return result, args.Error(1)
	}

	return nil, args.Error(1)
}
