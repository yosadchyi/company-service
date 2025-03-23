package postgresrepo

import (
	"company-service/internal/auth"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*auth.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, username, password FROM users WHERE username=$1", username)
	var u auth.User
	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil
}
