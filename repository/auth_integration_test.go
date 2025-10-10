package repository

import (
	"context"
	"native-free-pollings/database"
	"native-free-pollings/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{
		Email:        "dwiki@gmail.com",
		PasswordHash: "hash123",
		Name:         "dwiki",
	}

	repo := NewAuth(db)

	err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal("Failed create user:", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expected := &models.User{
		Email:        "dwiki@gmail.com",
		PasswordHash: "hash123",
		Name:         "dwiki",
	}

	repo := NewAuth(db)
	user, err := repo.GetUserByEmail(ctx, "dwiki@gmail.com")
	if err != nil {
		t.Fatal("user not found:", err)
	}

	assert.NotNil(t, user)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, expected.PasswordHash, user.PasswordHash)
	assert.Equal(t, expected.Name, user.Name)
}
