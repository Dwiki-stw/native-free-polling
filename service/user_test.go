package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"native-free-pollings/helper"
	"native-free-pollings/mocks"
	"native-free-pollings/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetProfile(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(repo *mocks.UserRepositoryMock)
		id         int64
		wantErr    string
	}{
		{
			name: "user not found",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("GetByID", mock.Anything, int64(1)).
					Return(nil, sql.ErrNoRows)
			},
			id:      1,
			wantErr: "NOT_FOUND",
		},
		{
			name: "db error",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("GetByID", mock.Anything, int64(2)).
					Return(nil, errors.New("db error"))
			},
			id:      2,
			wantErr: "INTERNAL_ERROR",
		},
		{
			name: "success",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("GetByID", mock.Anything, int64(3)).
					Return(&models.User{
						ID:        3,
						Name:      "test user",
						Email:     "test@example.com",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			id:      3,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepositoryMock)
			tt.setupMocks(repo)

			svc := NewUserService(repo, mocks.MockHasher{ShouldFail: false})
			resp, err := svc.GetProfile(context.Background(), tt.id)

			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.id, resp.ID)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(repo *mocks.UserRepositoryMock)
		user       *models.User
		wantErr    string
	}{
		{
			name: "user not found",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(sql.ErrNoRows)
			},
			user: &models.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: "test123",
				Name:         "test user",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: "NOT_FOUND",
		},
		{
			name: "db error",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(errors.New("db error"))
			},
			user: &models.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: "test123",
				Name:         "test user",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name: "user not found (get by id)",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil)
				repo.On("GetByID", mock.Anything, int64(1)).
					Return(nil, sql.ErrNoRows)

			},
			user: &models.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: "test123",
				Name:         "test user",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: "NOT_FOUND",
		},
		{
			name: "db error (get by id)",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil)
				repo.On("GetByID", mock.Anything, int64(2)).
					Return(nil, errors.New("db error"))
			},
			user: &models.User{
				ID:           2,
				Email:        "test@example.com",
				PasswordHash: "test123",
				Name:         "test user",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name: "success",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil)
				repo.On("GetByID", mock.Anything, int64(3)).
					Return(&models.User{
						ID:        3,
						Name:      "test user",
						Email:     "test@example.com",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			user: &models.User{
				ID:           3,
				Email:        "test@example.com",
				PasswordHash: "test123",
				Name:         "test user",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		repo := new(mocks.UserRepositoryMock)
		tt.setupMocks(repo)

		svc := NewUserService(repo, mocks.MockHasher{ShouldFail: false})
		resp, err := svc.UpdateProfile(context.Background(), tt.user)

		if tt.wantErr == "" {
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.user.ID, resp.ID)
		} else {
			assert.Nil(t, resp)
			assert.NotNil(t, err)
			assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
		}
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	tests := []struct {
		name       string
		id         int64
		password   string
		setupMocks func(repo *mocks.UserRepositoryMock)
		hasher     mocks.MockHasher
		wantErr    string
	}{
		{
			name:       "invalid input",
			id:         0,
			password:   "",
			setupMocks: func(repo *mocks.UserRepositoryMock) {},
			hasher:     mocks.MockHasher{ShouldFail: false},
			wantErr:    "BAD_REQUEST",
		},
		{
			name:       "hash password failed",
			id:         1,
			password:   strings.Repeat("a", 100),
			setupMocks: func(repo *mocks.UserRepositoryMock) {},
			hasher:     mocks.MockHasher{ShouldFail: true},
			wantErr:    "HASH_FAILED",
		},
		{
			name:     "user not found",
			id:       2,
			password: "hashed123",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("UpdatePassword", mock.Anything, int64(2), "hashed123").Return(sql.ErrNoRows)
			},
			hasher:  mocks.MockHasher{ShouldFail: false},
			wantErr: "NOT_FOUND",
		},
		{
			name:     "db error",
			id:       3,
			password: "hashed123",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("UpdatePassword", mock.Anything, int64(3), "hashed123").Return(errors.New("db error"))
			},
			hasher:  mocks.MockHasher{ShouldFail: false},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:     "success",
			id:       4,
			password: "hashed123",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("UpdatePassword", mock.Anything, int64(4), "hashed123").Return(nil)
			},
			hasher:  mocks.MockHasher{ShouldFail: false},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		repo := new(mocks.UserRepositoryMock)
		tt.setupMocks(repo)

		svc := NewUserService(repo, tt.hasher)
		err := svc.ChangePassword(context.Background(), tt.id, tt.password)

		if tt.wantErr == "" {
			assert.NoError(t, err)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
		}
	}
}
