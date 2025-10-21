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

func insertDummyVote(t *testing.T, db *sql.DB, vote *models.Vote) {
	query := `
		INSERT INTO votes(option_id, device_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.QueryRowContext(ctx, query, vote.OptionID, vote.DeviceHash).Scan(&vote.ID, &vote.CreatedAt)
	if err != nil {
		t.Fatalf("insert vote failed: %v", err)
	}
}

func insertDummyUserVote(t *testing.T, db *sql.DB, userID, voteID int64) {
	query := `
		INSERT INTO user_votes (user_id, vote_id)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, userID, voteID)
	if err != nil {
		t.Fatalf("insert vote failed: %v", err)
	}
}

func TestCreateVote(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vote := &models.Vote{
		OptionID:   1,
		DeviceHash: "device hash test",
	}

	repo := NewVote(db)
	err := repo.Create(ctx, db, vote)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, vote.ID)
}

func TestCreateUserVote(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	conf := Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vote := &models.Vote{
		OptionID:   1,
		DeviceHash: "device hash test",
	}

	insertDummyVote(t, db, vote)

	userVote := &models.UserVote{
		UserID: 1,
		VoteID: vote.ID,
	}

	repo := NewVote(db)
	err := repo.CreateUserVote(ctx, db, userVote.UserID, userVote.VoteID)

	assert.NoError(t, err)
}

func TestGetVotesByOptionID(t *testing.T) {
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

	insertDummyPolling(t, db, poll)

	option := &models.PollOption{
		PollID:   poll.ID,
		Label:    "Rust",
		Position: 1,
	}

	insertDummyOption(t, db, option)

	devices := []string{"device 1", "device 2", "device 3"}

	for _, d := range devices {
		vote := &models.Vote{
			OptionID:   option.ID,
			DeviceHash: d,
		}
		insertDummyVote(t, db, vote)
	}

	repo := NewVote(db)
	votes, err := repo.GetByOptionID(ctx, db, option.ID)

	assert.NoError(t, err)
	for i, d := range devices {
		assert.Equal(t, d, votes[i].DeviceHash)
	}
}

func TestGetVotesByPollID(t *testing.T) {
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

	insertDummyPolling(t, db, poll)

	option := &models.PollOption{
		PollID:   poll.ID,
		Label:    "Rust",
		Position: 1,
	}

	insertDummyOption(t, db, option)

	devices := []string{"device 1", "device 2", "device 3"}

	for _, d := range devices {
		vote := &models.Vote{
			OptionID:   option.ID,
			DeviceHash: d,
		}
		insertDummyVote(t, db, vote)
	}

	repo := NewVote(db)
	votes, err := repo.GetByPollID(ctx, db, poll.ID)

	assert.NoError(t, err)
	for i, d := range devices {
		assert.Equal(t, d, votes[i].DeviceHash)
	}

}

func TestHasDeviceVoted(t *testing.T) {
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

	insertDummyPolling(t, db, poll)

	option := &models.PollOption{
		PollID:   poll.ID,
		Label:    "Rust",
		Position: 1,
	}

	insertDummyOption(t, db, option)

	vote := &models.Vote{
		OptionID:   option.ID,
		DeviceHash: "device test",
	}
	insertDummyVote(t, db, vote)

	repo := NewVote(db)
	result, err := repo.HasDeviceVoted(ctx, db, vote.DeviceHash, poll.ID)

	assert.NoError(t, err)
	assert.Equal(t, true, result)
}

func TestHasUserVoted(t *testing.T) {
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

	insertDummyPolling(t, db, poll)

	option := &models.PollOption{
		PollID:   poll.ID,
		Label:    "Rust",
		Position: 1,
	}

	insertDummyOption(t, db, option)

	vote := &models.Vote{
		OptionID:   option.ID,
		DeviceHash: "device test",
	}
	insertDummyVote(t, db, vote)

	insertDummyUserVote(t, db, 3, vote.ID)

	repo := NewVote(db)
	result, err := repo.HasUserVoted(ctx, db, poll.ID, 3)

	assert.NoError(t, err)
	assert.Equal(t, true, result)
}
