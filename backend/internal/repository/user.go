package repository

import (
	"context"
	"database/sql"

	"github.com/osmanmertacar/elci/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context) (domain.User, error)
	Get(ctx context.Context, id int64) (domain.User, error)
}

type sqliteUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) Create(ctx context.Context) (domain.User, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO users DEFAULT VALUES`)
	if err != nil {
		return domain.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}
	return r.Get(ctx, id)
}

func (r *sqliteUserRepository) Get(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, `SELECT id, created_at FROM users WHERE id = ?`, id).
		Scan(&u.ID, &u.CreatedAt)
	return u, err
}
