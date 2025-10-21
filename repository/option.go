package repository

import (
	"context"
	"database/sql"
	"fmt"
	"native-free-pollings/domain"
	"native-free-pollings/models"
)

type option struct {
	DB *sql.DB
}

func NewOption(db *sql.DB) domain.OptionRepository {
	return &option{DB: db}
}

func (o *option) Create(ctx context.Context, db domain.DB, option *models.PollOption) error {
	query := `
		INSERT INTO poll_options (poll_id, label, position)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := db.QueryRowContext(ctx, query, option.PollID, option.Label, option.Position).Scan(&option.ID)
	if err != nil {
		return fmt.Errorf("insert option failed: %w", err)
	}

	return nil
}

func (o *option) Delete(ctx context.Context, db domain.DB, id int64) error {
	query := `
		DELETE FROM poll_options WHERE id  = $1
	`
	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete option failed: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (o *option) Update(ctx context.Context, db domain.DB, option *models.PollOption) error {
	query := `
		UPDATE poll_options
		SET label = $1,
			position = $2
		WHERE id = $3
	`
	result, err := db.ExecContext(ctx, query, option.Label, option.Position, option.ID)
	if err != nil {
		return fmt.Errorf("update option failed: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (o *option) GetByPollID(ctx context.Context, db domain.DB, id int64) ([]models.PollOption, error) {
	var options []models.PollOption

	query := `
		SELECT id, poll_id, label, position
		FROM poll_options
		WHERE poll_id = $1
		ORDER BY position
	`

	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query get options error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var opt models.PollOption
		err := rows.Scan(&opt.ID, &opt.PollID, &opt.Label, &opt.Position)
		if err != nil {
			return nil, fmt.Errorf("scan option failed: %w", err)
		}
		options = append(options, opt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows option iteration failed: %w", err)
	}

	return options, nil
}
