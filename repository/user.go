package repository

import (
	"context"
	"database/sql"
	"fmt"
	"native-free-pollings/domain"
	"native-free-pollings/models"
	"time"
)

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{DB: db}
}

func (u *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
	SELECT id, email, name, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	row := u.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
	UPDATE users
	SET name = $1,
	email = $2,
	updated_at = $3
	WHERE id = $4
	`
	result, err := u.DB.ExecContext(ctx, query,
		user.Name,
		user.Email,
		time.Now(),
		user.ID,
	)

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (u *userRepository) UpdatePassword(ctx context.Context, id int64, passwordHashed string) error {
	query := `
		UPDATE users
		SET password_hash = $1,
			updated_at = $2
		WHERE id = $3
	`

	result, err := u.DB.ExecContext(ctx, query,
		passwordHashed,
		time.Now(),
		id,
	)

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (u *userRepository) FindPollingsByID(ctx context.Context, id int64) ([]models.PollingSummary, error) {
	query := `
		SELECT p.id, p.title, p.status, count(v.id) 
		FROM polls p 
		LEFT JOIN poll_options po ON po.poll_id = p.id 
		LEFT JOIN votes v ON v.option_id = po.id 
		WHERE p.user_id = $1
		GROUP BY p.id, p.title, p.status
		ORDER BY p.created_at DESC
	`

	rows, err := u.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.PollingSummary
	for rows.Next() {
		var ps models.PollingSummary
		if err := rows.Scan(&ps.ID, &ps.Title, &ps.Status, &ps.TotalVotes); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("interation failed: %w", err)
	}

	return results, nil
}

func (u *userRepository) FindPollingsVotedByID(ctx context.Context, id int64) ([]models.PollingSummary, error) {
	query := `
		SELECT p.id, p.title, p.status, po."label"
		FROM polls p 
		JOIN poll_options po ON po.poll_id = p.id 
		JOIN votes v ON v.option_id = po.id 
		JOIN user_votes uv ON uv.vote_id = v.id 
		WHERE uv.user_id = $1
		ORDER BY p.created_at DESC
	`

	rows, err := u.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.PollingSummary
	for rows.Next() {
		var ps models.PollingSummary
		if err := rows.Scan(&ps.ID, &ps.Title, &ps.Status, &ps.UserVotedOption); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("interation failed: %w", err)
	}

	return results, nil
}
