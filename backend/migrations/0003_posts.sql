CREATE TABLE posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	default_caption TEXT NOT NULL DEFAULT '',
	default_media_kind TEXT,
	default_media_urls TEXT NOT NULL DEFAULT '[]',
	status TEXT NOT NULL DEFAULT 'draft',
	scheduled_at DATETIME,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_posts_status_scheduled_at ON posts (status, scheduled_at);

-- No connection_id: platform_connections has one row per (user_id,
-- platform), so a target's connection is always looked up by that pair.
CREATE TABLE post_targets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
	platform TEXT NOT NULL,
	caption_override TEXT,
	media_kind_override TEXT,
	media_urls_override TEXT,
	settings_json TEXT NOT NULL DEFAULT '{}',
	status TEXT NOT NULL DEFAULT 'pending',
	platform_post_id TEXT NOT NULL DEFAULT '',
	error TEXT NOT NULL DEFAULT '',
	published_at DATETIME,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_post_targets_post_id ON post_targets (post_id);
CREATE INDEX idx_post_targets_status ON post_targets (status);
