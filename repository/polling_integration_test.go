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

func insertDummyPolling(t *testing.T, db *sql.DB, poll *models.Polling) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO polls (user_id, title, description, status, starts_at, ends_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRowContext(ctx, query, poll.UserID, poll.Title, poll.Description, poll.Status, poll.StartsAt, poll.EndsAt).Scan(&poll.ID, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to insert dummy polling: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM polls WHERE id = $1`, poll.ID)
	})
}

func insertDummyOption(t *testing.T, db *sql.DB, opt *models.PollOption) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO poll_options (poll_id, label, position)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := db.QueryRowContext(ctx, query, opt.PollID, opt.Label, opt.Position).Scan(&opt.ID)
	if err != nil {
		t.Fatalf("failed to insert dummy option: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM poll_options WHERE id = $1`, opt.ID)
	})
}

func TestCreatePolling(t *testing.T) {
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
		Title:       "Title test",
		Description: "Description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	repo := NewPolling(db)
	err := repo.Create(ctx, db, poll)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, poll.ID)

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM polls WHERE id = $1`, poll.ID)
	})
}

func TestGetByIDPolling(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expectedPoll := &models.Polling{
		UserID:      1,
		Title:       "Title test",
		Description: "Description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	insertDummyPolling(t, db, expectedPoll)

	repo := NewPolling(db)
	poll, err := repo.GetByID(ctx, db, expectedPoll.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPoll.ID, poll.ID)
	assert.Equal(t, expectedPoll.Title, poll.Title)
}

func TestUpdatePolling(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dummyPoll := &models.Polling{
		UserID:      1,
		Title:       "Title test",
		Description: "Description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	insertDummyPolling(t, db, dummyPoll)

	updatePoll := &models.Polling{
		ID:          dummyPoll.ID,
		UserID:      1,
		Title:       "Title update",
		Description: "Description update",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	repo := NewPolling(db)
	err := repo.Update(ctx, db, updatePoll)

	assert.NoError(t, err)

	poll, err := repo.GetByID(ctx, db, updatePoll.ID)
	assert.NoError(t, err)
	assert.Equal(t, poll.Title, updatePoll.Title)
	assert.Equal(t, poll.Description, updatePoll.Description)
}

func TestGetResultByIDPolling(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dummyPoll := &models.Polling{
		UserID:      1,
		Title:       "Title test",
		Description: "Description test",
		Status:      "active",
		StartsAt:    time.Now(),
		EndsAt:      time.Now(),
	}

	insertDummyPolling(t, db, dummyPoll)

	labels := []string{"go", "kotlin", "javascript", "java"}
	for i, l := range labels {
		opt := &models.PollOption{
			PollID:   dummyPoll.ID,
			Label:    l,
			Position: i + 1,
		}
		insertDummyOption(t, db, opt)
	}

	repo := NewPolling(db)
	result, err := repo.GetResultsByID(ctx, db, dummyPoll.ID)

	assert.NoError(t, err)
	for i := range labels {
		assert.Equal(t, labels[i], result[i].OptionLabel)
	}
}
