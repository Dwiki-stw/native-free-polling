package repository

import (
	"context"
	"database/sql"
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
	_, err := u.DB.ExecContext(ctx, query,
		user.Name,
		user.Email,
		time.Now(),
		user.ID,
	)

	return err
}

func (u *userRepository) UpdatePassword(ctx context.Context, id int64, passwordHashed string) error {
	query := `
		UPDATE users
		SET password_hash = $1,
			updated_at = $2
		WHERE id = $3
	`

	_, err := u.DB.ExecContext(ctx, query,
		passwordHashed,
		time.Now(),
		id,
	)

	return err
}
