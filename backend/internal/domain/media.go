package domain

import "time"

type MediaKind string

const (
	MediaImage MediaKind = "image"
	MediaVideo MediaKind = "video"
)

type MediaAsset struct {
	ID              int64
	UserID          int64
	Kind            MediaKind
	StorageKey      string
	PublicURL       string
	ContentType     string
	SizeBytes       int64
	DurationSeconds *int
	CreatedAt       time.Time
}
