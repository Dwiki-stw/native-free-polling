package repository

import (
	"context"
	"database/sql"
	"native-free-pollings/domain"
	"native-free-pollings/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupAuthMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, domain.AuthRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	repo := NewAuth(db)

	return db, mock, repo
}

func TestAuthRepository_CreateUser(t *testing.T) {
	db, mock, repo := setupAuthMockDB(t)
	defer db.Close()

	user := &models.User{
		Email:        "test@example.com",
		PasswordHash: "test123",
		Name:         "test user",
	}

	expectedId := int64(1)
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(expectedId, expectedCreatedAt, expectedUpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users (email, password_hash, name)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
	`)).
		WithArgs(user.Email, user.PasswordHash, user.Name).
		WillReturnRows(rows)

	err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, expectedId, user.ID)
	assert.Equal(t, expectedCreatedAt, user.CreatedAt)
	assert.Equal(t, expectedUpdatedAt, user.UpdatedAt)
}

func TestAuthRepository_GetByEmail(t *testing.T) {
	db, mock, repo := setupAuthMockDB(t)
	defer db.Close()

	expedtedUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "test123",
		Name:         "test user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "name", "created_at", "updated_at"}).
		AddRow(expedtedUser.ID, expedtedUser.Email, expedtedUser.PasswordHash, expedtedUser.Name, expedtedUser.CreatedAt, expedtedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, password_hash, name, created_at, updated_at
        FROM users
        WHERE email = $1
		LIMIT 1 
	`)).
		WithArgs(expedtedUser.Email).
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail(context.Background(), expedtedUser.Email)

	assert.NoError(t, err)
	assert.Equal(t, expedtedUser.ID, user.ID)
	assert.Equal(t, expedtedUser.Email, user.Email)
	assert.Equal(t, expedtedUser.Name, user.Name)
}
