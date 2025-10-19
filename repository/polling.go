package repository

import (
	"context"
	"database/sql"
	"fmt"
	"native-free-pollings/domain"
	"native-free-pollings/models"
	"time"
)

type polling struct {
	DB *sql.DB
}

func NewPolling(db *sql.DB) domain.PollRepository {
	return &polling{DB: db}
}

func (p *polling) Create(ctx context.Context, db domain.DB, poll *models.Polling) error {
	query := `
		INSERT INTO polls (user_id, title, description, status, starts_at, ends_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRowContext(ctx, query, poll.UserID, poll.Title, poll.Description, poll.Status, poll.StartsAt, poll.EndsAt).
		Scan(&poll.ID, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert polling failed: %w", err)
	}

	return nil
}

func (p *polling) Delete(ctx context.Context, db domain.DB, id int64) error {
	query := `
		DELETE polls WHERE id = $1
	`
	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete poll failed: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *polling) Update(ctx context.Context, db domain.DB, poll *models.Polling) error {
	query := `
		UPDATE polls
		SET title = $1,
			description = $2,
			status = $3,
			stars_at = $4,
			ends_at = $5,
			updated_at = $6
		WHERE id = $7
	`
	result, err := db.ExecContext(ctx, query, poll.Title, poll.Description, poll.Status, poll.StartsAt, poll.EndsAt, time.Now())

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (p *polling) GetByID(ctx context.Context, db domain.DB, id int64) (*models.Polling, error) {
	var poll models.Polling

	query := `
		SELECT p.id, p.user_id, p.title, p.description
			   p.status, p.starts_at, p.ends_at, p.created_at,
			   p.updated.at, u.name as creator_name, u.email as creator_email 
		FROM polls p
		JOIN users u ON u.id = p.user_id
		WHERE id = $1
	`
	err := db.QueryRowContext(ctx, query, id).
		Scan(&poll.ID, &poll.UserID, &poll.Title, &poll.Description, &poll.Status, &poll.StartsAt, &poll.EndsAt, &poll.CreatedAt, &poll.UpdatedAt, &poll.CreatorName, &poll.CreatorEmail)
	if err != nil {
		return nil, fmt.Errorf("get polling failed: %w", err)
	}

	query = `
		SELECT o.id, o.poll_id, o.label, o.position
		FROM poll_options o
		WHERE o.poll_id = $1
		ORDER BY o.position
	`
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query get options error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var opt models.PollOption
		if err := rows.Scan(&opt.ID, &opt.PollID, &opt.Label, &opt.Position); err != nil {
			return nil, fmt.Errorf("scan options failed: %w", err)
		}
		poll.Options = append(poll.Options, opt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows option iteration failed: %w", err)
	}

	return &poll, nil
}

func (p *polling) GetResultsByID(ctx context.Context, db domain.DB, id int64) ([]models.VoteResult, error) {
	query := `
		SELECT o.id as option_id, o.label as option label, COUNT(v.id) as votes
		FROM poll_options o
		LEFT JOIN votes v ON v.option_id = o.id
		WHERE o.poll_id = $1
		GROUP BY o.id, o.label, o.position
		ORDER BY o.position
	`
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.VoteResult
	for rows.Next() {
		var rest models.VoteResult
		if err := rows.Scan(); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results, rest)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("interation failed: %w", err)
	}

	return results, nil
}
