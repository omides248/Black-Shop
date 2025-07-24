package postgresql

import (
	"black-shop-service/internal/domain/identity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type userRepo struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewUserRepository(conn *pgx.Conn, logger *zap.Logger) (identity.UserRepository, error) {
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT,
			email TEXT UNIQUE,
			password_hash TEXT
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return &userRepo{
		conn:   conn,
		logger: logger.Named("postgres_repo"),
	}, nil
}

func (r *userRepo) Save(ctx context.Context, user *identity.User) error {
	r.logger.Debug("saving a new user", zap.String("email", user.Email))

	query := "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id"

	var userID string
	err := r.conn.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			r.logger.Warn("attempted to save user with duplicate email", zap.String("email", user.Email))
			return identity.ErrEmailAlreadyExists
		}
		r.logger.Error("failed to save user", zap.Error(err))
		return err
	}

	user.ID = identity.UserID(userID)
	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*identity.User, error) {
	r.logger.Info("finding user by email", zap.String("email", email))

	var user identity.User
	query := "SELECT id, name, email, password_hash FROM users WHERE email = $1"

	err := r.conn.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("user not found by email", zap.String("email", email))
			return nil, identity.ErrUserNotFound
		}
		r.logger.Error("failed to find user by email", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByID(ctx context.Context, id identity.UserID) (*identity.User, error) {
	r.logger.Info("finding user by id", zap.String("user_id", string(id)))

	var user identity.User
	query := "SELECT id, name, email, password_hash FROM users WHERE id = $1"

	err := r.conn.QueryRow(ctx, query, string(id)).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("user not found by id", zap.String("user_id", string(id)))
			return nil, identity.ErrUserNotFound
		}
		r.logger.Error("failed to find user by id", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
