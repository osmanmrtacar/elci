package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/domain"
)

type OAuthSessionRepository interface {
	Create(ctx context.Context, s domain.OAuthSession) (domain.OAuthSession, error)
	// Consume fetches a session by state and deletes it, so each state can
	// only ever be redeemed once.
	Consume(ctx context.Context, state string) (domain.OAuthSession, error)
}

type sqliteOAuthSessionRepository struct {
	db *sql.DB
}

func NewOAuthSessionRepository(db *sql.DB) OAuthSessionRepository {
	return &sqliteOAuthSessionRepository{db: db}
}

func (r *sqliteOAuthSessionRepository) Create(ctx context.Context, s domain.OAuthSession) (domain.OAuthSession, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO oauth_sessions (state, code_verifier, platform, user_id, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`, s.State, s.CodeVerifier, s.Platform, s.UserID, s.ExpiresAt)
	if err != nil {
		return domain.OAuthSession{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.OAuthSession{}, err
	}
	s.ID = id
	return s, nil
}

func (r *sqliteOAuthSessionRepository) Consume(ctx context.Context, state string) (domain.OAuthSession, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.OAuthSession{}, err
	}
	defer tx.Rollback()

	var s domain.OAuthSession
	err = tx.QueryRowContext(ctx, `
		SELECT id, state, code_verifier, platform, user_id, created_at, expires_at
		FROM oauth_sessions WHERE state = ?
	`, state).Scan(&s.ID, &s.State, &s.CodeVerifier, &s.Platform, &s.UserID, &s.CreatedAt, &s.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.OAuthSession{}, ErrNotFound
	}
	if err != nil {
		return domain.OAuthSession{}, err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM oauth_sessions WHERE id = ?`, s.ID); err != nil {
		return domain.OAuthSession{}, err
	}

	if time.Now().After(s.ExpiresAt) {
		return domain.OAuthSession{}, errors.New("oauth session expired")
	}

	return s, tx.Commit()
}
