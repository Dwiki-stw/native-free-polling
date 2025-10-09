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

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, domain.UserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("faild to open sqlmock: %v", err)
	}

	repo := NewUserRepository(db)
	return db, mock, repo
}

func TestUserRepository_GetById(t *testing.T) {
	_, mock, repo := setupMockDB(t)

	expectedUser := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Name, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`)).
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	user, err := repo.GetByID(context.Background(), expectedUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	user := &models.User{
		ID:    1,
		Email: "updated@example.com",
		Name:  "updated User",
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET name = $1,
			email = $2,
			updated_at = $3
		WHERE id = $4
	`)).
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	userID := int64(1)
	hashedPassword := "hashedPassword"

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET password_hash = $1,
			updated_at = $2
		WHERE id = $3
	`)).
		WithArgs(hashedPassword, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdatePassword(context.Background(), userID, hashedPassword)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
