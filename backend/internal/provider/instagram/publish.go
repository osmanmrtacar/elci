package instagram

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

type containerResponse struct {
	ID string `json:"id"`
}

type containerStatusResponse struct {
	StatusCode string `json:"status_code"`
}

// Publish creates a media container, waits for Instagram to finish
// processing it (required for video/Reels, usually instant for images),
// then publishes it. By the time this returns, the post is live — there is
// nothing left for PollStatus to wait on.
func (p *Provider) Publish(ctx context.Context, token provider.Token, content provider.Content, _ map[string]any) (provider.PublishRef, error) {
	form := url.Values{
		"access_token": {token.AccessToken},
		"caption":      {content.Caption},
	}
	if content.MediaKind == provider.MediaVideo {
		form.Set("video_url", content.MediaURLs[0])
		form.Set("media_type", "REELS")
	} else {
		form.Set("image_url", content.MediaURLs[0])
	}

	var container containerResponse
	if err := postGraph(ctx, "/me/media", form, &container); err != nil {
		return provider.PublishRef{}, fmt.Errorf("instagram: create container: %w", err)
	}

	if err := waitUntilFinished(ctx, token.AccessToken, container.ID); err != nil {
		return provider.PublishRef{}, err
	}

	var published containerResponse
	err := postGraph(ctx, "/me/media_publish", url.Values{
		"access_token": {token.AccessToken},
		"creation_id":  {container.ID},
	}, &published)
	if err != nil {
		return provider.PublishRef{}, fmt.Errorf("instagram: publish container: %w", err)
	}

	return provider.PublishRef{ID: published.ID}, nil
}

func waitUntilFinished(ctx context.Context, accessToken, containerID string) error {
	const maxAttempts = 10
	const pollInterval = 30 * time.Second

	for attempt := 0; attempt < maxAttempts; attempt++ {
		var status containerStatusResponse
		err := getGraph(ctx, "/"+containerID, url.Values{
			"fields":       {"status_code"},
			"access_token": {accessToken},
		}, &status)
		if err != nil {
			return fmt.Errorf("instagram: check container status: %w", err)
		}

		switch status.StatusCode {
		case "FINISHED":
			return nil
		case "ERROR", "EXPIRED":
			return fmt.Errorf("instagram: media processing %s for container %s", status.StatusCode, containerID)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}

	return fmt.Errorf("instagram: media processing timed out for container %s", containerID)
}
