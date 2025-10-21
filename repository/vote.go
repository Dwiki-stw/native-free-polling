package repository

import (
	"context"
	"database/sql"
	"fmt"
	"native-free-pollings/domain"
	"native-free-pollings/models"
)

type vote struct {
	DB *sql.DB
}

func NewVote(db *sql.DB) domain.VoteRepository {
	return &vote{DB: db}
}

func (v *vote) Create(ctx context.Context, db domain.DB, vote *models.Vote) error {
	query := `
		INSERT INTO votes(option_id, device_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := db.QueryRowContext(ctx, query, vote.OptionID, vote.DeviceHash).Scan(&vote.ID, &vote.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert vote failed: %w", err)
	}

	return nil
}

func (v *vote) CreateUserVote(ctx context.Context, db domain.DB, userID int64, voteID int64) error {
	query := `
		INSERT INTO user_votes (user_id, vote_id)
		VALUES ($1, $2)
	`
	_, err := db.ExecContext(ctx, query, userID, voteID)
	if err != nil {
		return fmt.Errorf("insert user vote failed: %w", err)
	}

	return nil
}

func (v *vote) GetByOptionID(ctx context.Context, db domain.DB, optionID int64) ([]models.Vote, error) {
	query := `
		SELECT v.id, v.option_id, v.device_hash, v.created_at
		FROM votes v
		WHERE v.option_id = $1
	`

	rows, err := db.QueryContext(ctx, query, optionID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		var vt models.Vote
		err := rows.Scan(&vt.ID, &vt.OptionID, &vt.DeviceHash, &vt.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan vote failed: %w", err)
		}
		votes = append(votes, vt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("interation vote rows failed: %w", err)
	}

	return votes, nil
}

func (v *vote) GetByPollID(ctx context.Context, db domain.DB, pollID int64) ([]models.Vote, error) {
	query := `
		SELECT v.id, v.option_id, v.device_hash, v.created_at
		FROM votes v
		JOIN poll_options o ON o.id = v.option_id
		WHERE  o.poll_id = $1
	`

	rows, err := db.QueryContext(ctx, query, pollID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		var vt models.Vote
		err := rows.Scan(&vt.ID, &vt.OptionID, &vt.DeviceHash, &vt.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan vote failed: %w", err)
		}
		votes = append(votes, vt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("interation vote rows failed: %w", err)
	}

	return votes, nil
}

func (v *vote) HasDeviceVoted(ctx context.Context, db domain.DB, deviceHash string, pollID int64) (bool, error) {
	var exist bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM votes v
			JOIN poll_options o ON o.id = v.option_id
			WHERE v.device_hash = $1 and o.poll_id = $2
		)
	`

	err := db.QueryRowContext(ctx, query, deviceHash, pollID).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("query error : %w", err)
	}

	return exist, nil
}

func (v *vote) HasUserVoted(ctx context.Context, db domain.DB, pollID int64, userID int64) (bool, error) {
	var exist bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_votes uv
			JOIN votes v ON v.id = uv.vote_id
			JOIN poll_options o ON o.id = v.option_id
			WHERE uv.user_id = $1 and o.poll_id = $2
		)
	`

	err := db.QueryRowContext(ctx, query, userID, pollID).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("query error : %w", err)
	}

	return exist, nil
}
