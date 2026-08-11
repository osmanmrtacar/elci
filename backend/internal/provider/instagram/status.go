package instagram

import (
	"context"
	"net/url"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

// PollStatus is a formality here too: Publish already waited for processing
// and called media_publish before returning, so this just confirms the
// resulting media still exists.
func (p *Provider) PollStatus(ctx context.Context, token provider.Token, ref provider.PublishRef) (provider.Status, error) {
	var res struct {
		ID string `json:"id"`
	}
	err := getGraph(ctx, "/"+ref.ID, url.Values{
		"fields":       {"id"},
		"access_token": {token.AccessToken},
	}, &res)
	if err != nil {
		return provider.Status{State: "failed", Message: err.Error()}, nil
	}
	return provider.Status{State: "published"}, nil
}
