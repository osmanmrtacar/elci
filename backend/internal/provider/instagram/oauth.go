package instagram

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

const shortLivedTokenURL = "https://api.instagram.com/oauth/access_token"
const longLivedExchangeURL = "https://graph.instagram.com/access_token"
const refreshURL = "https://graph.instagram.com/refresh_access_token"

type shortLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type longLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (p *Provider) Exchange(ctx context.Context, code, _ string) (provider.Token, provider.AccountInfo, error) {
	var shortLived shortLivedTokenResponse
	err := doGraph(ctx, http.MethodPost, shortLivedTokenURL, url.Values{
		"client_id":     {p.cfg.AppID},
		"client_secret": {p.cfg.AppSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {p.cfg.RedirectURI},
		"code":          {code},
	}, &shortLived)
	if err != nil {
		return provider.Token{}, provider.AccountInfo{}, err
	}

	tok, err := p.exchangeForLongLived(ctx, shortLived.AccessToken)
	if err != nil {
		return provider.Token{}, provider.AccountInfo{}, err
	}

	info, err := p.AccountInfo(ctx, tok)
	if err != nil {
		return provider.Token{}, provider.AccountInfo{}, err
	}

	return tok, info, nil
}

func (p *Provider) exchangeForLongLived(ctx context.Context, shortLivedAccessToken string) (provider.Token, error) {
	var res longLivedTokenResponse
	err := doGraph(ctx, http.MethodGet, longLivedExchangeURL, url.Values{
		"grant_type":    {"ig_exchange_token"},
		"client_secret": {p.cfg.AppSecret},
		"access_token":  {shortLivedAccessToken},
	}, &res)
	if err != nil {
		return provider.Token{}, err
	}

	return provider.Token{
		AccessToken: res.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(res.ExpiresIn) * time.Second),
	}, nil
}

// RefreshToken renews a long-lived token; Instagram has no refresh_token
// concept — the same access token is extended in place.
func (p *Provider) RefreshToken(ctx context.Context, token provider.Token) (provider.Token, error) {
	var res longLivedTokenResponse
	err := doGraph(ctx, http.MethodGet, refreshURL, url.Values{
		"grant_type":   {"ig_refresh_token"},
		"access_token": {token.AccessToken},
	}, &res)
	if err != nil {
		return provider.Token{}, err
	}

	return provider.Token{
		AccessToken: res.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(res.ExpiresIn) * time.Second),
	}, nil
}
