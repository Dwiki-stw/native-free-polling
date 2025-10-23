package service

import (
	"context"
	"database/sql"
	"errors"
	"native-free-pollings/dto"
	"native-free-pollings/helper"
	"native-free-pollings/mocks"
	"native-free-pollings/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type BundleMockPoll struct {
	PollRepo *mocks.PollRepositoryMock
	OptRepo  *mocks.OptionRepositoryMock
	VoteRepo *mocks.VoteRepostoryMock
}

func TestPollingService_CreatePolling(t *testing.T) {
	tests := []struct {
		name       string
		req        *dto.CreatePollingRequest
		setupMocks func(repo *BundleMockPoll)
		setupDB    func(db *sql.DB, mock sqlmock.Sqlmock)
		wantErr    string
	}{
		{
			name:       "error begin tx",
			req:        &dto.CreatePollingRequest{Title: "test create poll", Description: "test description create poll", Options: []string{"Go"}},
			setupMocks: func(repo *BundleMockPoll) {},
			setupDB:    func(db *sql.DB, mock sqlmock.Sqlmock) {},
			wantErr:    "INTERNAL_ERROR",
		},
		{
			name: "error create polling",
			req:  &dto.CreatePollingRequest{Title: "test create poll", Description: "test description create poll", Options: []string{"Go"}},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(errors.New("error create polling"))
			},
			setupDB: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name: "error create option",
			req:  &dto.CreatePollingRequest{Title: "test create poll", Description: "test description create poll", Options: []string{"Go"}},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(errors.New("error create polling"))

				repo.OptRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.PollOption")).
					Return(errors.New("failed to save option"))
			},
			setupDB: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name: "error tx commit",
			req:  &dto.CreatePollingRequest{Title: "test create poll", Description: "test description create poll", Options: []string{"Go"}},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Run(func(args mock.Arguments) {
						poll := args.Get(2).(*models.Polling)
						poll.ID = 1
						poll.CreatedAt = time.Now()
						poll.UpdatedAt = time.Now()
					}).
					Return(nil)

				repo.OptRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.PollOption")).
					Return(nil)
			},
			setupDB: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name: "success",
			req:  &dto.CreatePollingRequest{Title: "test create poll", Description: "test description create poll", Options: []string{"Go"}},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Run(func(args mock.Arguments) {
						poll := args.Get(2).(*models.Polling)
						poll.ID = 1
						poll.CreatedAt = time.Now()
						poll.UpdatedAt = time.Now()
					}).
					Return(nil)

				repo.OptRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.PollOption")).
					Return(nil)
			},
			setupDB: func(db *sql.DB, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupDB(db, mock)

			bundleMock := &BundleMockPoll{
				PollRepo: new(mocks.PollRepositoryMock),
				OptRepo:  new(mocks.OptionRepositoryMock),
				VoteRepo: new(mocks.VoteRepostoryMock),
			}

			tt.setupMocks(bundleMock)

			svc := NewPolling(db, bundleMock.PollRepo, bundleMock.OptRepo, bundleMock.VoteRepo)
			creator := dto.CreatorInfo{ID: 1, Name: "user test", Email: "test@example.com"}

			resp, err := svc.CreatePolling(context.Background(), tt.req, creator)

			if tt.wantErr != "" {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Title, resp.Title)
			}
		})
	}
}

func TestPollingService_UpdatePolling(t *testing.T) {
	tests := []struct {
		name       string
		req        *dto.UpdatePollingRequest
		creator    dto.CreatorInfo
		setupMocks func(repo *BundleMockPoll)
		setupDB    func(mock sqlmock.Sqlmock)
		wantErr    string
	}{
		{
			name:    "error get polling",
			req:     &dto.UpdatePollingRequest{},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil, errors.New("failed get polling"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:    "error not creator polling",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description"},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      2,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "FORBIDDEN_ERROR",
		},
		{
			name:    "error get options",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description"},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil, errors.New("failed get options"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:    "error begin tx",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description"},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{
						{ID: 1},
						{ID: 2},
						{ID: 3},
					}, nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:    "error update polling",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description"},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{
						{ID: 1},
						{ID: 2},
						{ID: 3},
					}, nil)
				repo.PollRepo.On("Update", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(errors.New("failed update polling"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name:    "error update option",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description", Options: []dto.Option{{ID: 1}}},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{
						{ID: 1},
						{ID: 2},
						{ID: 3},
					}, nil)
				repo.PollRepo.On("Update", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(nil)
				repo.OptRepo.On("Update", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.PollOption")).
					Return(errors.New("failed update option"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name:    "error create option",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description", Options: []dto.Option{{ID: 4}}},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{
						{ID: 1},
						{ID: 2},
						{ID: 3},
					}, nil)
				repo.PollRepo.On("Update", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(nil)
				repo.OptRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.PollOption")).
					Return(errors.New("failed create option"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name:    "error delete option",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description"},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{
						{ID: 1},
						{ID: 2},
						{ID: 3},
					}, nil)
				repo.PollRepo.On("Update", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(nil)
				repo.OptRepo.On("Delete", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("int64")).
					Return(errors.New("failed delete option"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name:    "error commit tx",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description"},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{}, nil)
				repo.PollRepo.On("Update", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:    "success",
			req:     &dto.UpdatePollingRequest{ID: 1, Title: "Test title", Description: "Test description"},
			creator: dto.CreatorInfo{ID: 1, Name: "test user", Email: "test@example.com"},
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						ID:          1,
						UserID:      1,
						Title:       "Test title",
						Description: "Test description",
					}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{}, nil)
				repo.PollRepo.On("Update", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Polling")).
					Return(nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupDB(mock)

			bundleMock := &BundleMockPoll{
				PollRepo: new(mocks.PollRepositoryMock),
				OptRepo:  new(mocks.OptionRepositoryMock),
				VoteRepo: new(mocks.VoteRepostoryMock),
			}

			tt.setupMocks(bundleMock)

			svc := NewPolling(db, bundleMock.PollRepo, bundleMock.OptRepo, bundleMock.VoteRepo)

			resp, err := svc.UpdatePolling(context.Background(), tt.req, tt.creator)

			if tt.wantErr != "" {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.ID, resp.ID)
			}
		})
	}
}

func TestPollingService_VoteOptionPolling(t *testing.T) {
	tests := []struct {
		name       string
		userID     int64
		pollID     int64
		optionID   int64
		setupMocks func(repo *BundleMockPoll)
		setupDB    func(mock sqlmock.Sqlmock)
		wantErr    string
	}{
		{
			name:     "error has user voted",
			userID:   1,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasUserVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(false, errors.New("failed get user voted"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:     "error user exist",
			userID:   1,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasUserVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(true, nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "ALREADY_VOTED",
		},
		{
			name:     "error has device voted",
			userID:   0,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasDeviceVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
					Return(false, errors.New("failed get device voted"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:     "error device exist",
			userID:   0,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasDeviceVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
					Return(true, nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "ALREADY_VOTED",
		},
		{
			name:     "error begin tx",
			userID:   0,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasDeviceVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
					Return(false, nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:     "error get polling",
			userID:   0,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasDeviceVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
					Return(false, nil)
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil, errors.New("failed get polling"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name:     "error polling status not active",
			userID:   0,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasDeviceVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
					Return(false, nil)
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						Status: "draft",
					}, nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "BAD_REQUEST",
		},
		{
			name:     "error create vote",
			userID:   0,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasDeviceVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
					Return(false, nil)
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						Status: "active",
					}, nil)
				repo.VoteRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Vote")).
					Return(errors.New("failed create vote"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name:     "error create user vote",
			userID:   1,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasUserVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(false, nil)
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						Status: "active",
					}, nil)
				repo.VoteRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Vote")).
					Return(nil)
				repo.VoteRepo.On("CreateUserVote", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(errors.New("failed create user vote"))
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			wantErr: "DB_ERROR",
		},
		{
			name:     "error commit tx",
			userID:   1,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasUserVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(false, nil)
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						Status: "active",
					}, nil)
				repo.VoteRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Vote")).
					Return(nil)
				repo.VoteRepo.On("CreateUserVote", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(errors.New("failed commit"))
			},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name:     "success",
			userID:   1,
			pollID:   1,
			optionID: 1,
			setupMocks: func(repo *BundleMockPoll) {
				repo.VoteRepo.On("HasUserVoted", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(false, nil)
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{
						Status: "active",
					}, nil)
				repo.VoteRepo.On("Create", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("*models.Vote")).
					Return(nil)
				repo.VoteRepo.On("CreateUserVote", mock.Anything, mock.IsType(&sql.Tx{}), mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).
					Return(nil)
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.setupDB(mock)

			bundleMock := &BundleMockPoll{
				PollRepo: new(mocks.PollRepositoryMock),
				OptRepo:  new(mocks.OptionRepositoryMock),
				VoteRepo: new(mocks.VoteRepostoryMock),
			}

			tt.setupMocks(bundleMock)

			svc := NewPolling(db, bundleMock.PollRepo, bundleMock.OptRepo, bundleMock.VoteRepo)

			err := svc.VoteOptionPolling(context.Background(), tt.userID, tt.pollID, tt.optionID, "example")

			if tt.wantErr != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPollingService_DeletePolling(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(repo *BundleMockPoll)
		wantErr    string
	}{
		{
			name: "error get polling",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil, errors.New("failed get polling"))
			},
			wantErr: "INTERNAL_ERROR",
		},
		{
			name: "error not creator",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{ID: 1, UserID: 0}, nil)
			},
			wantErr: "FORBIDDEN_ERROR",
		},
		{
			name: "error delete polling",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{ID: 1, UserID: 1}, nil)
				repo.PollRepo.On("Delete", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(errors.New("failed delete polling"))
			},
			wantErr: "DB_ERROR",
		},
		{
			name: "success",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{ID: 1, UserID: 1}, nil)
				repo.PollRepo.On("Delete", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil)
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, _ := sqlmock.New()
			defer db.Close()

			bundleMock := &BundleMockPoll{
				PollRepo: new(mocks.PollRepositoryMock),
				OptRepo:  new(mocks.OptionRepositoryMock),
				VoteRepo: new(mocks.VoteRepostoryMock),
			}

			tt.setupMocks(bundleMock)

			svc := NewPolling(db, bundleMock.PollRepo, bundleMock.OptRepo, bundleMock.VoteRepo)

			err := svc.DeletePolling(context.Background(), 1, 1)

			if tt.wantErr != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPollingService_GetDetailPolling(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(repo *BundleMockPoll)
		wantErr    string
	}{
		{
			name: "error get polling",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil, errors.New("failed get polling"))
			},
			wantErr: "DB_ERROR",
		},
		{
			name: "error get options",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{ID: 1}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil, errors.New("failed get options"))
			},
			wantErr: "DB_ERROR",
		},
		{
			name: "success",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(&models.Polling{ID: 1}, nil)
				repo.OptRepo.On("GetByPollID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.PollOption{{}}, nil)
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, _ := sqlmock.New()
			defer db.Close()

			bundleMock := &BundleMockPoll{
				PollRepo: new(mocks.PollRepositoryMock),
				OptRepo:  new(mocks.OptionRepositoryMock),
				VoteRepo: new(mocks.VoteRepostoryMock),
			}

			tt.setupMocks(bundleMock)

			svc := NewPolling(db, bundleMock.PollRepo, bundleMock.OptRepo, bundleMock.VoteRepo)

			resp, err := svc.GetDetailPolling(context.Background(), 1)

			if tt.wantErr != "" {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
func TestPollingService_GetPollingResult(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(repo *BundleMockPoll)
		wantErr    string
	}{
		{
			name: "error get polling result",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetResultsByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return(nil, errors.New("failed get polling"))
			},
			wantErr: "DB_ERROR",
		},
		{
			name: "success",
			setupMocks: func(repo *BundleMockPoll) {
				repo.PollRepo.On("GetResultsByID", mock.Anything, mock.IsType(&sql.DB{}), mock.AnythingOfType("int64")).
					Return([]models.VoteResult{
						{OptionID: 1, OptionLabel: "go", Votes: 10},
						{OptionID: 2, OptionLabel: "go", Votes: 10},
						{OptionID: 3, OptionLabel: "go", Votes: 10},
					}, nil)
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, _ := sqlmock.New()
			defer db.Close()

			bundleMock := &BundleMockPoll{
				PollRepo: new(mocks.PollRepositoryMock),
				OptRepo:  new(mocks.OptionRepositoryMock),
				VoteRepo: new(mocks.VoteRepostoryMock),
			}

			tt.setupMocks(bundleMock)

			svc := NewPolling(db, bundleMock.PollRepo, bundleMock.OptRepo, bundleMock.VoteRepo)

			resp, err := svc.GetPollingResult(context.Background(), 1)

			if tt.wantErr != "" {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
