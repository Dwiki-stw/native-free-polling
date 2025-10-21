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

func getOption(db *sql.DB, id int64) *models.PollOption {
	query := `
		SELECT * FROM poll_options WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var option models.PollOption
	err := db.QueryRowContext(ctx, query, id).Scan(&option.ID, &option.PollID, &option.Label, &option.Position, &option.CreatedAt)
	if err != nil {
		return &option
	}

	return &option
}

func TestCreateOption(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poll := &models.Polling{
		UserID:      1,
		Title:       "polling test",
		Description: "polling description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	insertDummyPolling(t, db, poll)

	option := &models.PollOption{
		PollID:   poll.ID,
		Label:    "Rust",
		Position: 1,
	}

	repo := NewOption(db)
	err := repo.Create(ctx, db, option)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, option.ID)
}

func TestUpdateOption(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poll := &models.Polling{
		UserID:      1,
		Title:       "polling test",
		Description: "polling description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	insertDummyPolling(t, db, poll)

	oldOption := &models.PollOption{
		PollID:   poll.ID,
		Label:    "Rust",
		Position: 1,
	}

	insertDummyOption(t, db, oldOption)

	updateOption := &models.PollOption{
		ID:       oldOption.ID,
		PollID:   poll.ID,
		Label:    "Go",
		Position: 2,
	}

	repo := NewOption(db)
	err := repo.Update(ctx, db, updateOption)

	assert.NoError(t, err)

	option := getOption(db, updateOption.ID)

	assert.Equal(t, updateOption.Label, option.Label)
	assert.Equal(t, updateOption.Position, option.Position)
}

func TestDeleteOption(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poll := &models.Polling{
		UserID:      1,
		Title:       "polling test",
		Description: "polling description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	insertDummyPolling(t, db, poll)

	option := &models.PollOption{
		PollID:   poll.ID,
		Label:    "Rust",
		Position: 1,
	}

	insertDummyOption(t, db, option)

	repo := NewOption(db)
	err := repo.Delete(ctx, db, option.ID)

	assert.NoError(t, err)

	result := getOption(db, option.ID)

	assert.NotEqual(t, option.Label, result.Label)
	assert.NotEqual(t, option.Position, result.Position)
}

func TestGetOptionByPollID(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poll := &models.Polling{
		UserID:      1,
		Title:       "polling test",
		Description: "polling description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	insertDummyPolling(t, db, poll)

	labels := []string{"go", "kotlin", "java", "rust"}

	for i, l := range labels {
		option := &models.PollOption{
			PollID:   poll.ID,
			Label:    l,
			Position: i + 1,
		}
		insertDummyOption(t, db, option)
	}

	repo := NewOption(db)
	options, err := repo.GetByPollID(ctx, db, poll.ID)

	assert.NoError(t, err)
	for i, l := range labels {
		assert.Equal(t, l, options[i].Label)
		assert.Equal(t, i+1, options[i].Position)
	}
}
