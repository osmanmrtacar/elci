package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/domain"
)

var ErrNotFound = errors.New("not found")

type ConnectionRepository interface {
	FindByPlatformUserID(ctx context.Context, platform domain.Platform, platformUserID string) (domain.PlatformConnection, error)
	Get(ctx context.Context, userID int64, platform domain.Platform) (domain.PlatformConnection, error)
	ListByUser(ctx context.Context, userID int64) ([]domain.PlatformConnection, error)
	Upsert(ctx context.Context, c domain.PlatformConnection) (domain.PlatformConnection, error)
	Deactivate(ctx context.Context, userID int64, platform domain.Platform) error
}

type sqliteConnectionRepository struct {
	db *sql.DB
}

func NewConnectionRepository(db *sql.DB) ConnectionRepository {
	return &sqliteConnectionRepository{db: db}
}

const connectionColumns = `
	id, user_id, platform, platform_user_id, username, display_name, avatar_url,
	access_token, refresh_token, token_expires_at, scope, is_active, connected_at, last_used_at
`

func scanConnection(row interface{ Scan(...any) error }) (domain.PlatformConnection, error) {
	var c domain.PlatformConnection
	var expiresAt sql.NullTime
	err := row.Scan(
		&c.ID, &c.UserID, &c.Platform, &c.PlatformUserID, &c.Username, &c.DisplayName, &c.AvatarURL,
		&c.AccessToken, &c.RefreshToken, &expiresAt, &c.Scope, &c.IsActive, &c.ConnectedAt, &c.LastUsedAt,
	)
	if err != nil {
		return domain.PlatformConnection{}, err
	}
	c.TokenExpiresAt = expiresAt.Time
	return c, nil
}

func (r *sqliteConnectionRepository) FindByPlatformUserID(ctx context.Context, platform domain.Platform, platformUserID string) (domain.PlatformConnection, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+connectionColumns+` FROM platform_connections WHERE platform = ? AND platform_user_id = ?`, platform, platformUserID)
	c, err := scanConnection(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PlatformConnection{}, ErrNotFound
	}
	return c, err
}

func (r *sqliteConnectionRepository) Get(ctx context.Context, userID int64, platform domain.Platform) (domain.PlatformConnection, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+connectionColumns+` FROM platform_connections WHERE user_id = ? AND platform = ?`, userID, platform)
	c, err := scanConnection(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PlatformConnection{}, ErrNotFound
	}
	return c, err
}

func (r *sqliteConnectionRepository) ListByUser(ctx context.Context, userID int64) ([]domain.PlatformConnection, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+connectionColumns+` FROM platform_connections WHERE user_id = ? ORDER BY connected_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.PlatformConnection
	for rows.Next() {
		c, err := scanConnection(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *sqliteConnectionRepository) Upsert(ctx context.Context, c domain.PlatformConnection) (domain.PlatformConnection, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO platform_connections (
			user_id, platform, platform_user_id, username, display_name, avatar_url,
			access_token, refresh_token, token_expires_at, scope, is_active, last_used_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?)
		ON CONFLICT (user_id, platform) DO UPDATE SET
			platform_user_id = excluded.platform_user_id,
			username = excluded.username,
			display_name = excluded.display_name,
			avatar_url = excluded.avatar_url,
			access_token = excluded.access_token,
			refresh_token = excluded.refresh_token,
			token_expires_at = excluded.token_expires_at,
			scope = excluded.scope,
			is_active = 1,
			last_used_at = excluded.last_used_at
	`,
		c.UserID, c.Platform, c.PlatformUserID, c.Username, c.DisplayName, c.AvatarURL,
		c.AccessToken, c.RefreshToken, c.TokenExpiresAt, c.Scope, time.Now(),
	)
	if err != nil {
		return domain.PlatformConnection{}, err
	}
	return r.Get(ctx, c.UserID, c.Platform)
}

func (r *sqliteConnectionRepository) Deactivate(ctx context.Context, userID int64, platform domain.Platform) error {
	_, err := r.db.ExecContext(ctx, `UPDATE platform_connections SET is_active = 0 WHERE user_id = ? AND platform = ?`, userID, platform)
	return err
}
