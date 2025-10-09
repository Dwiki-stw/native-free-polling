package repository

import (
	"context"
	"database/sql"
	"native-free-pollings/database"
	"native-free-pollings/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func insertDummy(t *testing.T, db *sql.DB, email, name, password string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64

	query := `
		INSERT INTO users(email, name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	//add user
	err := db.QueryRowContext(ctx, query, email, name, password).Scan(&id)

	if err != nil {
		t.Fatalf("failed to insert dummy user: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	})

	return id
}

func TestGetUserById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expectedEmail := "test@example.com"
	expectedName := "test user"
	id := insertDummy(t, db, expectedEmail, expectedName, "password123")

	repo := NewUserRepository(db)

	user, err := repo.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmail, user.Email)
	assert.Equal(t, expectedName, user.Name)
}

func TestUpdateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := insertDummy(t, db, "old@example.com", "old user", "password123")

	expectedUser := &models.User{
		ID:    id,
		Email: "updated@example.com",
		Name:  "updated User",
	}

	repo := NewUserRepository(db)
	err := repo.Update(ctx, expectedUser)
	assert.NoError(t, err)

	user, err := repo.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)
}

func TestUpdatePasswordUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := insertDummy(t, db, "test@example.com", "test user", "oldpassword123")

	expedtedPassword := "updated123"
	repo := NewUserRepository(db)

	err := repo.UpdatePassword(ctx, id, expedtedPassword)
	assert.NoError(t, err)

	var updatePassword string
	err = db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE id = $1`, id).
		Scan(&updatePassword)
	assert.NoError(t, err)
	assert.Equal(t, expedtedPassword, updatePassword)
}
