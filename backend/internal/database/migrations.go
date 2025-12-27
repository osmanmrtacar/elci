package database

import (
	"fmt"
	"log"
)

// RunMigrations runs all database migrations
func (db *DB) RunMigrations() error {
	log.Println("Running database migrations...")

	migrations := []string{
		createUsersTable,
		createTokensTable,
		createPostsTable,
		createPlatformConnectionsTable,
		createOAuthSessionsTable,
		createIndexes,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tiktok_user_id TEXT UNIQUE NOT NULL,
    username TEXT,
    display_name TEXT,
    avatar_url TEXT,
    platform TEXT DEFAULT 'tiktok',
    platform_user_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createTokensTable = `
CREATE TABLE IF NOT EXISTS tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_type TEXT DEFAULT 'Bearer',
    expires_at TIMESTAMP NOT NULL,
    scope TEXT,
    platform TEXT DEFAULT 'tiktok',
    code_verifier TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`

const createPostsTable = `
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    tiktok_post_id TEXT,
    platform_post_id TEXT,
    video_url TEXT NOT NULL,
    caption TEXT,
    platform TEXT DEFAULT 'tiktok',
    media_type TEXT DEFAULT 'video',
    status TEXT DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_platform_user ON users(platform, platform_user_id);
CREATE INDEX IF NOT EXISTS idx_tokens_user_platform ON tokens(user_id, platform);
CREATE INDEX IF NOT EXISTS idx_posts_user_platform ON posts(user_id, platform);
CREATE INDEX IF NOT EXISTS idx_posts_platform_status ON posts(platform, status);
CREATE INDEX IF NOT EXISTS idx_platform_connections_user ON platform_connections(user_id);
CREATE INDEX IF NOT EXISTS idx_platform_connections_platform ON platform_connections(platform);
CREATE INDEX IF NOT EXISTS idx_oauth_sessions_state ON oauth_sessions(state);
CREATE INDEX IF NOT EXISTS idx_oauth_sessions_expires ON oauth_sessions(expires_at);
`

// Create platform_connections table for tracking connected platforms per user
const createPlatformConnectionsTable = `
CREATE TABLE IF NOT EXISTS platform_connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    platform TEXT NOT NULL,
    platform_user_id TEXT NOT NULL,
    username TEXT,
    display_name TEXT,
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, platform)
);
`

// Create oauth_sessions table for PKCE state management
const createOAuthSessionsTable = `
CREATE TABLE IF NOT EXISTS oauth_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    state TEXT UNIQUE NOT NULL,
    code_verifier TEXT,
    platform TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);
`
