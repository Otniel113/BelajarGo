package repositories

import (
	"context"
	"database/sql"
	"errors"

	"cashier/models"
)

type AuthRepository interface {
	Register(ctx context.Context, u *models.User) error
	FindByUsernameOrEmail(ctx context.Context, identity string) (*models.User, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(ctx context.Context, u *models.User) error {
	query := "INSERT INTO users (username, email, password, status) VALUES ($1, $2, $3, $4) RETURNING id"
	return r.db.QueryRowContext(ctx, query, u.Username, u.Email, u.Password, u.Status).Scan(&u.ID)
}

func (r *authRepository) FindByUsernameOrEmail(ctx context.Context, identity string) (*models.User, error) {
	query := "SELECT id, username, email, password, status FROM users WHERE username = $1 OR email = $1"
	u := &models.User{}
	err := r.db.QueryRowContext(ctx, query, identity).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *authRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, username, email, password, status FROM users WHERE id = $1"
	u := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}
