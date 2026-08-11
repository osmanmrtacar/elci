package instagram

import (
	"context"
	"net/url"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

type meResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AccountType string `json:"account_type"`
}

func (p *Provider) AccountInfo(ctx context.Context, token provider.Token) (provider.AccountInfo, error) {
	var res meResponse
	query := url.Values{
		"fields":       {"id,username,account_type"},
		"access_token": {token.AccessToken},
	}
	if err := getGraph(ctx, "/me", query, &res); err != nil {
		return provider.AccountInfo{}, err
	}

	return provider.AccountInfo{
		PlatformUserID: res.ID,
		Username:       res.Username,
		DisplayName:    res.Username,
		Extra: map[string]any{
			"account_type": res.AccountType,
		},
	}, nil
}
