package x

import (
	"context"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

// PollStatus is a formality for X: POST /2/tweets already published the post
// synchronously by the time Publish returned, so this just confirms it's
// still there rather than waiting on any further processing.
func (p *Provider) PollStatus(ctx context.Context, token provider.Token, ref provider.PublishRef) (provider.Status, error) {
	var res tweetResponse
	err := bearerJSON(ctx, "GET", token.AccessToken, "/2/tweets/"+ref.ID, nil, &res)
	if err != nil {
		return provider.Status{State: "failed", Message: err.Error()}, nil
	}
	return provider.Status{State: "published"}, nil
}
