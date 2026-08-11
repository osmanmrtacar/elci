package tiktok

import (
	"context"
	"net/url"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

const tokenURL = "https://open.tiktokapis.com/v2/oauth/token/"

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	OpenID       string `json:"open_id"`
	TokenType    string `json:"token_type"`
}

func (p *Provider) Exchange(ctx context.Context, code, _ string) (provider.Token, provider.AccountInfo, error) {
	form := url.Values{
		"client_key":    {p.cfg.ClientKey},
		"client_secret": {p.cfg.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {p.cfg.RedirectURI},
	}

	var res tokenResponse
	if err := postOAuthForm(ctx, tokenURL, form, &res); err != nil {
		return provider.Token{}, provider.AccountInfo{}, err
	}

	tok := provider.Token{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(res.ExpiresIn) * time.Second),
		Scope:        res.Scope,
	}

	info, err := p.AccountInfo(ctx, tok)
	if err != nil {
		return provider.Token{}, provider.AccountInfo{}, err
	}

	return tok, info, nil
}

func (p *Provider) RefreshToken(ctx context.Context, token provider.Token) (provider.Token, error) {
	form := url.Values{
		"client_key":    {p.cfg.ClientKey},
		"client_secret": {p.cfg.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {token.RefreshToken},
	}

	var res tokenResponse
	if err := postOAuthForm(ctx, tokenURL, form, &res); err != nil {
		return provider.Token{}, err
	}

	return provider.Token{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(res.ExpiresIn) * time.Second),
		Scope:        res.Scope,
	}, nil
}
