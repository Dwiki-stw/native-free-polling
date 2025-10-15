package handler

import (
	"context"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"native-free-pollings/mocks"
	"native-free-pollings/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerGetProfile(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         any
		setupMocks func(svc *mocks.UserServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "method no allowed",
			method:     http.MethodPost,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			wantCode:   http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name:       "invalid user id",
			method:     http.MethodGet,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			id:         "not int64",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "invalid user id",
		},
		{
			name:   "service return error",
			method: http.MethodGet,
			id:     int64(1),
			setupMocks: func(svc *mocks.UserServiceMock) {
				svc.On("GetProfile", mock.Anything, int64(1)).
					Return(nil, helper.NewAppError("DB_ERROR", "db failed", nil))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `"code":"DB_ERROR"`,
		},
		{
			name:   "success",
			method: http.MethodGet,
			id:     int64(2),
			setupMocks: func(svc *mocks.UserServiceMock) {
				svc.On("GetProfile", mock.Anything, int64(2)).
					Return(&dto.ProfileResponse{
						ID:    2,
						Name:  "test user",
						Email: "test@example.com",
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `"email":"test@example.com"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.UserServiceMock)
			tt.setupMocks(svc)

			h := &UserHandler{Service: svc}

			req := httptest.NewRequest(tt.method, "/profile", nil)
			if tt.id != nil {
				req = req.WithContext(context.WithValue(req.Context(), "userID", tt.id))
			}
			rr := httptest.NewRecorder()

			h.GetProfile(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerUpdateProfile(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         any
		body       string
		setupMocks func(svc *mocks.UserServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "method no allowed",
			method:     http.MethodGet,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			wantCode:   http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name:       "invalid user id",
			method:     http.MethodPost,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			id:         "not int64",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "invalid user id",
		},
		{
			name:       "invalid request",
			method:     http.MethodPost,
			id:         int64(1),
			body:       `{"email": "a@mail.com", "name":}`,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name:       "bad request",
			method:     http.MethodPost,
			id:         int64(1),
			body:       `{"email":"", "name":""}`,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "name and email cannot both be empty",
		},
		{
			name:   "service return error",
			method: http.MethodPost,
			id:     int64(2),
			body:   `{"email":"test@example.com", "name":"test user"}`,
			setupMocks: func(svc *mocks.UserServiceMock) {
				svc.On("UpdateProfile", mock.Anything, &models.User{ID: 2, Email: "test@example.com", Name: "test user"}).
					Return(nil, helper.NewAppError("DB_ERROR", "db failed", nil))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `"code":"DB_ERROR"`,
		},
		{
			name:   "success",
			method: http.MethodPost,
			id:     int64(3),
			body:   `{"email":"test@example.com", "name":"test user"}`,
			setupMocks: func(svc *mocks.UserServiceMock) {
				svc.On("UpdateProfile", mock.Anything, &models.User{ID: 3, Email: "test@example.com", Name: "test user"}).
					Return(&dto.ProfileResponse{
						ID:    3,
						Name:  "test user",
						Email: "test@example.com",
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `"email":"test@example.com"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.UserServiceMock)
			tt.setupMocks(svc)

			h := &UserHandler{Service: svc}

			req := httptest.NewRequest(tt.method, "/profile", strings.NewReader(tt.body))
			if tt.id != nil {
				req = req.WithContext(context.WithValue(req.Context(), "userID", tt.id))
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.UpdateProfile(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerChangePassword(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         any
		body       string
		setupMocks func(svc *mocks.UserServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "method no allowed",
			method:     http.MethodGet,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			wantCode:   http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name:       "invalid user id",
			method:     http.MethodPost,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			id:         "not int64",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "invalid user id",
		},
		{
			name:       "invalid request",
			method:     http.MethodPost,
			id:         int64(1),
			body:       `{"password":}`,
			setupMocks: func(svc *mocks.UserServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name:   "service return error",
			method: http.MethodPost,
			id:     int64(2),
			body:   `{"password":"test123"}`,
			setupMocks: func(svc *mocks.UserServiceMock) {
				svc.On("ChangePassword", mock.Anything, int64(2), "test123").
					Return(helper.NewAppError("DB_ERROR", "db failed", nil))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `"code":"DB_ERROR"`,
		},
		{
			name:   "success",
			method: http.MethodPost,
			id:     int64(3),
			body:   `{"password":"test123"}`,
			setupMocks: func(svc *mocks.UserServiceMock) {
				svc.On("ChangePassword", mock.Anything, int64(3), "test123").
					Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: `"message":"password updated successfully"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.UserServiceMock)
			tt.setupMocks(svc)

			h := &UserHandler{Service: svc}

			req := httptest.NewRequest(tt.method, "/users/me/password", strings.NewReader(tt.body))
			if tt.id != nil {
				req = req.WithContext(context.WithValue(req.Context(), "userID", tt.id))
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.ChangePassword(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}
