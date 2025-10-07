package handler

import (
	"native-free-pollings/dto"
	"native-free-pollings/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerRegister(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		setupMocks func(svc *mocks.AuthServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:   "method not allowed",
			method: http.MethodGet,
			body:   ``,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				// tidak ada expectation
			},
			wantCode: http.StatusMethodNotAllowed,
			wantBody: "method not allowed\n",
		},
		{
			name:   "invalid request payload",
			method: http.MethodPost,
			body:   `{"email": "a@mail.com", "password": 123}`, // pass bukan string
			setupMocks: func(svc *mocks.AuthServiceMock) {
				// tidak ada expectation
			},
			wantCode: http.StatusBadRequest,
			wantBody: "invalid request payload\n",
		},
		{
			name:   "service returns error",
			method: http.MethodPost,
			body:   `{"email":"a@mail.com","password":"123","name":"John"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("Register", mock.Anything, &dto.RegisterRequest{
					Email: "a@mail.com",
					Pass:  "123",
					Name:  "John",
				}).Return(nil, assert.AnError)
			},
			wantCode: http.StatusBadRequest,
			wantBody: assert.AnError.Error() + "\n",
		},
		{
			name:   "success",
			method: http.MethodPost,
			body:   `{"email":"a@mail.com","password":"123","name":"John"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("Register", mock.Anything, &dto.RegisterRequest{
					Email: "a@mail.com",
					Pass:  "123",
					Name:  "John",
				}).Return(&dto.RegisterResponse{
					ID:    1,
					Email: "a@mail.com",
					Name:  "John",
				}, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: `"email":"a@mail.com"`, // cukup cek sebagian JSON
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.AuthServiceMock)
			tt.setupMocks(svc)

			h := &AuthHandler{Service: svc}

			req := httptest.NewRequest(tt.method, "/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.Register(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerLogin(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		setupMocks func(svc *mocks.AuthServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "method not allowed",
			body:       ``,
			method:     http.MethodGet,
			setupMocks: func(svc *mocks.AuthServiceMock) {},
			wantCode:   http.StatusMethodNotAllowed,
			wantBody:   "method not allowed\n",
		},
		{
			name:       "invalid request payload",
			body:       ``,
			method:     http.MethodPost,
			setupMocks: func(svc *mocks.AuthServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid request payload\n",
		},
		{
			name:   "service returns error",
			body:   `{"email": "a@mail.com", "password":"123"}`,
			method: http.MethodPost,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("Login", mock.Anything, &dto.LoginRequest{
					Email:    "a@mail.com",
					Password: "123",
				}).Return(nil, assert.AnError)
			},
			wantCode: http.StatusBadRequest,
			wantBody: assert.AnError.Error() + "\n",
		},
		{
			name:   "success",
			body:   `{"email": "a@mail.com", "password":"123"}`,
			method: http.MethodPost,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("Login", mock.Anything, &dto.LoginRequest{
					Email:    "a@mail.com",
					Password: "123",
				}).Return(&dto.LoginResponse{
					ID:    1,
					Name:  "John",
					Token: "newtoken",
				}, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: `"id":1`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.AuthServiceMock)
			tt.setupMocks(svc)

			h := &AuthHandler{Service: svc}

			req := httptest.NewRequest(tt.method, "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.Login(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}
