package handler

import (
	"context"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"native-free-pollings/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerCreatePolling(t *testing.T) {
	tests := []struct {
		name       string
		creator    any
		body       string
		setupMocks func(svc *mocks.PollServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "failed get creator information",
			creator:    "",
			body:       ``,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusUnauthorized,
			wantBody:   "invalid user information",
		},
		{
			name:       "invalid request body",
			creator:    &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body:       `{invalid json}`,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name:       "error validate request",
			creator:    &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body:       `{}`,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "payload validation failed",
		},
		{
			name:    "CreatePolling service error",
			creator: &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body:    `{"title": "title polling", "description": "description polling", "status": "active", "starts_at": "2025-10-25T21:52:00+07:00", "ends_at": "2025-10-25T21:52:00+07:00", "options": ["go", "kotlin"]}`,
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("CreatePolling", mock.Anything, mock.AnythingOfType("*dto.CreatePollingRequest"), mock.AnythingOfType("dto.CreatorInfo")).
					Return(nil, helper.NewAppError("BAD_REQUEST", assert.AnError.Error(), assert.AnError))
			},
			wantCode: http.StatusBadRequest,
			wantBody: assert.AnError.Error(),
		},
		{
			name:    "success",
			creator: &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body:    `{"title": "title polling", "description": "description polling", "status": "active", "starts_at": "2025-10-25T21:52:00+07:00", "ends_at": "2025-10-25T21:52:00+07:00", "options": ["go", "kotlin"]}`,
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("CreatePolling", mock.Anything, mock.AnythingOfType("*dto.CreatePollingRequest"), mock.AnythingOfType("dto.CreatorInfo")).
					Return(&dto.PollingResponse{ID: 1}, nil)
			},
			wantCode: http.StatusCreated,
			wantBody: "created polling successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.PollServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest("CREATE", "/pollings/1", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), helper.AuthKey, tt.creator)
			req = req.WithContext(ctx)

			h := &Polling{Service: svc}
			h.CreatePolling(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerUpdatePolling(t *testing.T) {
	tests := []struct {
		name       string
		creator    any
		body       string
		setupMocks func(svc *mocks.PollServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "failed get creator information",
			creator:    "",
			body:       ``,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusUnauthorized,
			wantBody:   "invalid user information",
		},
		{
			name:       "invalid request body",
			creator:    &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body:       `{invalid json}`,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid request payload",
		},
		{
			name:       "error validate request",
			creator:    &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body:       `{}`,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "payload validation failed",
		},
		{
			name:    "UpdatePolling service error",
			creator: &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body: `{
						"id": 1,
						"title": "Favorite Programming Language",
						"description": "Vote for the language you love most",
						"status": "active",
						"starts_at": "2025-10-26T09:00:00+07:00",
						"ends_at": "2025-10-30T17:00:00+07:00",
						"options": [
							{
								"id": 1,
								"label": "Go",
								"position": 1
							}
						]
					}`,
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("UpdatePolling", mock.Anything, mock.AnythingOfType("*dto.UpdatePollingRequest"), mock.AnythingOfType("dto.CreatorInfo")).
					Return(nil, helper.NewAppError("BAD_REQUEST", assert.AnError.Error(), assert.AnError))
			},
			wantCode: http.StatusBadRequest,
			wantBody: assert.AnError.Error(),
		},
		{
			name:    "success",
			creator: &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			body: `{
						"id": 1,
						"title": "Favorite Programming Language",
						"description": "Vote for the language you love most",
						"status": "active",
						"starts_at": "2025-10-26T09:00:00+07:00",
						"ends_at": "2025-10-30T17:00:00+07:00",
						"options": [
							{
							"id": 1,
							"label": "Go",
							"position": 1
							}
						]
					}`,
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("UpdatePolling", mock.Anything, mock.AnythingOfType("*dto.UpdatePollingRequest"), mock.AnythingOfType("dto.CreatorInfo")).
					Return(&dto.PollingResponse{ID: 1}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "updated polling successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.PollServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest("PATCH", "/pollings/1", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), helper.AuthKey, tt.creator)
			req = req.WithContext(ctx)

			h := &Polling{Service: svc}
			h.UpdatePolling(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerDeletePolling(t *testing.T) {
	tests := []struct {
		name       string
		creator    any
		path       string
		setupMocks func(svc *mocks.PollServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "failed get information creator",
			creator:    "",
			path:       "/pollings",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusUnauthorized,
			wantBody:   "invalid user information",
		},
		{
			name:       "invalid path",
			creator:    &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			path:       "/pollings",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid path id polling",
		},
		{
			name:       "invalid id polling",
			creator:    &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			path:       "/polling/",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid id polling",
		},
		{
			name:    "DeletePolling return error",
			creator: &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			path:    "/polling/1",
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("DeletePolling", mock.Anything, int64(1), int64(1)).
					Return(helper.NewAppError("NOT_FOUND", assert.AnError.Error(), assert.AnError))
			},
			wantCode: http.StatusNotFound,
			wantBody: assert.AnError.Error(),
		},
		{
			name:    "success",
			creator: &helper.AuthContext{UserID: 1, UserName: "test user", UserEmail: "test@example.com"},
			path:    "/polling/1",
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("DeletePolling", mock.Anything, int64(1), int64(1)).
					Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: "deleted polling successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.PollServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest("DELETE", tt.path, nil)
			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), helper.AuthKey, tt.creator)
			req = req.WithContext(ctx)

			h := &Polling{Service: svc}
			h.DeletePolling(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerGetDetailPolling(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		setupMocks func(svc *mocks.PollServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "invalid path",
			path:       "/pollings",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid path id polling",
		},
		{
			name:       "invalid id polling",
			path:       "/polling/",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid id polling",
		},
		{
			name: "GetDetailPolling return error",
			path: "/polling/1",
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("GetDetailPolling", mock.Anything, int64(1)).
					Return(nil, helper.NewAppError("NOT_FOUND", assert.AnError.Error(), assert.AnError))
			},
			wantCode: http.StatusNotFound,
			wantBody: assert.AnError.Error(),
		},
		{
			name: "success",
			path: "/polling/1",
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("GetDetailPolling", mock.Anything, int64(1)).
					Return(&dto.PollingResponse{}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "get detail polling successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.PollServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()

			h := &Polling{Service: svc}
			h.GetDetailPolling(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerVoteOptionPolling(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		setupMocks func(svc *mocks.PollServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "method invalid",
			method:     "GET",
			path:       "/pollings",
			body:       `{}`,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name:       "invalid path",
			method:     "POST",
			path:       "/pollings",
			body:       `{}`,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid path id polling",
		},
		{
			name:       "invalid id polling",
			method:     "POST",
			path:       "/pollings//vote",
			body:       `{}`,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid id polling",
		},
		{
			name:       "invalid request body",
			method:     "POST",
			path:       "/pollings//vote",
			body:       ``,
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid id polling",
		},
		{
			name:   "VoteOptionPolling return error",
			method: "POST",
			path:   "/pollings/1/vote",
			body:   `{"option_id": 1, "device_hash": "test device hash"}`,
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("VoteOptionPolling", mock.Anything, int64(0), int64(1), int64(1), "test device hash").
					Return(helper.NewAppError("NOT_FOUND", assert.AnError.Error(), assert.AnError))
			},
			wantCode: http.StatusNotFound,
			wantBody: assert.AnError.Error(),
		},
		{
			name:   "success",
			method: "POST",
			path:   "/polling/1/vote",
			body:   `{"option_id": 1, "device_hash": "test device hash"}`,
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("VoteOptionPolling", mock.Anything, int64(0), int64(1), int64(1), "test device hash").
					Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: "vote successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.PollServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), helper.AuthKey, "")
			req = req.WithContext(ctx)

			h := &Polling{Service: svc}
			h.VoteOptionPolling(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandlerGetPollingResult(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		setupMocks func(svc *mocks.PollServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "method invalid",
			method:     http.MethodPost,
			path:       "/",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusMethodNotAllowed,
			wantBody:   "method not allowed",
		},
		{
			name:       "invalid path",
			method:     http.MethodGet,
			path:       "/pollings",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid path id polling",
		},
		{
			name:       "invalid id polling",
			method:     http.MethodGet,
			path:       "/polling//result",
			setupMocks: func(svc *mocks.PollServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   "invalid id polling",
		},
		{
			name:   "GetPollingResult return error",
			method: http.MethodGet,
			path:   "/pollings/1/result",
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("GetPollingResult", mock.Anything, int64(1)).
					Return(nil, helper.NewAppError("NOT_FOUND", assert.AnError.Error(), assert.AnError))
			},
			wantCode: http.StatusNotFound,
			wantBody: assert.AnError.Error(),
		},
		{
			name:   "success",
			method: http.MethodGet,
			path:   "/pollings/1/result",
			setupMocks: func(svc *mocks.PollServiceMock) {
				svc.On("GetPollingResult", mock.Anything, int64(1)).
					Return(&dto.ResultPolling{}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "get polling result successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.PollServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			h := &Polling{Service: svc}
			h.GetPollingResult(rr, req)

			assert.Equal(t, tt.wantCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}
