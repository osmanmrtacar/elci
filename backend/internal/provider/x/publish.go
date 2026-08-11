package x

import (
	"context"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

type tweetResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// Publish uploads any media first (X has no PULL_FROM_URL equivalent, so the
// bytes are downloaded from our storage and re-uploaded), then posts the
// tweet. Unlike TikTok, this completes synchronously — there is nothing to
// poll afterward.
func (p *Provider) Publish(ctx context.Context, token provider.Token, content provider.Content, _ map[string]any) (provider.PublishRef, error) {
	mediaIDs := make([]string, 0, len(content.MediaURLs))
	for _, mediaURL := range content.MediaURLs {
		id, err := uploadMedia(ctx, token.AccessToken, mediaURL, content.MediaKind)
		if err != nil {
			return provider.PublishRef{}, err
		}
		mediaIDs = append(mediaIDs, id)
	}

	body := map[string]any{"text": content.Caption}
	if len(mediaIDs) > 0 {
		body["media"] = map[string]any{"media_ids": mediaIDs}
	}

	var res tweetResponse
	if err := bearerJSON(ctx, "POST", token.AccessToken, "/2/tweets", body, &res); err != nil {
		return provider.PublishRef{}, err
	}

	return provider.PublishRef{ID: res.Data.ID}, nil
}
