package mocks

import (
	"context"
	"native-free-pollings/dto"

	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.LoginResponse); ok {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *AuthServiceMock) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.RegisterResponse); ok {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}
