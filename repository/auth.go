package repository

import (
	"context"
	"database/sql"
	"native-free-pollings/domain"
	"native-free-pollings/models"
)

type auth struct {
	DB *sql.DB
}

func NewAuth(db *sql.DB) domain.AuthRepository {
	return &auth{DB: db}
}

func (a *auth) CreateUser(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (email, password_hash, name)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `
	return a.DB.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Name).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
}

func (a *auth) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, name, created_at, updated_at
        FROM users
        WHERE email = $1
		LIMIT 1
    `

	row := a.DB.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
