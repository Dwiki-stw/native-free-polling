package service

import (
	"context"
	"errors"
	"native-free-pollings/dto"
	"native-free-pollings/mocks"
	"native-free-pollings/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(repo *mocks.UserRepositoryMock)
		hasher     mocks.MockHasher
		req        *dto.RegisterRequest
		wantErr    string
	}{
		{
			name: "email alread exists",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "exist@mail.com").
					Return(&models.User{Email: "exist@mail.com"}, nil)
			},
			hasher:  mocks.MockHasher{ShouldFail: false},
			req:     &dto.RegisterRequest{Email: "exist@mail.com", Pass: "123", Name: "Dwiki"},
			wantErr: "EMAIL_EXIST",
		},
		{
			name: "hash password failed (simulate bcrypt error)",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "ok@mail.com").
					Return(nil, nil)
			},
			hasher:  mocks.MockHasher{ShouldFail: true},
			req:     &dto.RegisterRequest{Email: "ok@mail.com", Pass: strings.Repeat("a", 100), Name: "John"},
			wantErr: "HASH_FAILED",
		},
		{
			name: "db save failed",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "db@mail.com").
					Return(nil, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(errors.New("db error"))
			},
			hasher:  mocks.MockHasher{ShouldFail: false},
			req:     &dto.RegisterRequest{Email: "db@mail.com", Pass: "123", Name: "John"},
			wantErr: "DB_ERROR",
		},
		{
			name: "success",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "new@mail.com").
					Return(nil, nil)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil)
			},
			hasher:  mocks.MockHasher{ShouldFail: false},
			req:     &dto.RegisterRequest{Email: "new@mail.com", Pass: "123", Name: "John"},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.UserRepositoryMock)
			tt.setupMocks(repo)

			svc := NewAuthService(repo, []byte("test-secret"), tt.hasher)
			resp, err := svc.Register(context.Background(), tt.req)

			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Email, resp.Email)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, resp)
			}
		})
	}

}
