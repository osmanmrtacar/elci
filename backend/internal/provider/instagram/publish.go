package instagram

import (
	"context"
	"fmt"
	"net/url"
	"strings"
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
//
// Multiple image URLs are published as a carousel: each image becomes its
// own item container, then a parent CAROUSEL container references them by
// ID. A single image (or a video) skips straight to the plain container
// Instagram's API expects in that case.
func (p *Provider) Publish(ctx context.Context, token provider.Token, content provider.Content, _ map[string]any) (provider.PublishRef, error) {
	var containerID string
	var err error

	switch {
	case content.MediaKind == provider.MediaVideo:
		containerID, err = p.createContainer(ctx, token.AccessToken, url.Values{
			"video_url":  {content.MediaURLs[0]},
			"media_type": {"REELS"},
			"caption":    {content.Caption},
		})
	case len(content.MediaURLs) > 1:
		containerID, err = p.createCarouselContainer(ctx, token.AccessToken, content)
	default:
		containerID, err = p.createContainer(ctx, token.AccessToken, url.Values{
			"image_url": {content.MediaURLs[0]},
			"caption":   {content.Caption},
		})
	}
	if err != nil {
		return provider.PublishRef{}, err
	}

	if err := waitUntilFinished(ctx, token.AccessToken, containerID); err != nil {
		return provider.PublishRef{}, err
	}

	var published containerResponse
	err = postGraph(ctx, "/me/media_publish", url.Values{
		"access_token": {token.AccessToken},
		"creation_id":  {containerID},
	}, &published)
	if err != nil {
		return provider.PublishRef{}, fmt.Errorf("instagram: publish container: %w", err)
	}

	return provider.PublishRef{ID: published.ID}, nil
}

func (p *Provider) createContainer(ctx context.Context, accessToken string, form url.Values) (string, error) {
	form.Set("access_token", accessToken)

	var container containerResponse
	if err := postGraph(ctx, "/me/media", form, &container); err != nil {
		return "", fmt.Errorf("instagram: create container: %w", err)
	}
	return container.ID, nil
}

// createCarouselContainer creates one item container per image, waits for
// each to finish processing (Instagram requires children to be FINISHED
// before they can be attached to a carousel), then creates the parent
// CAROUSEL container referencing them.
func (p *Provider) createCarouselContainer(ctx context.Context, accessToken string, content provider.Content) (string, error) {
	childIDs := make([]string, 0, len(content.MediaURLs))
	for _, imageURL := range content.MediaURLs {
		id, err := p.createContainer(ctx, accessToken, url.Values{
			"image_url":        {imageURL},
			"is_carousel_item": {"true"},
		})
		if err != nil {
			return "", err
		}
		if err := waitUntilFinished(ctx, accessToken, id); err != nil {
			return "", err
		}
		childIDs = append(childIDs, id)
	}

	return p.createContainer(ctx, accessToken, url.Values{
		"media_type": {"CAROUSEL"},
		"children":   {strings.Join(childIDs, ",")},
		"caption":    {content.Caption},
	})
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
