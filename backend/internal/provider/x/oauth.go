package x

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

const tokenURL = apiBase + "/2/oauth2/token"

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (p *Provider) postTokenForm(ctx context.Context, form url.Values) (tokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return tokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	basic := base64.StdEncoding.EncodeToString([]byte(p.cfg.ClientID + ":" + p.cfg.ClientSecret))
	req.Header.Set("Authorization", "Basic "+basic)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokenResponse{}, err
	}

	if resp.StatusCode >= 300 {
		var oauthErr struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		_ = json.Unmarshal(body, &oauthErr)
		return tokenResponse{}, fmt.Errorf("x oauth: %s: %s", oauthErr.Error, oauthErr.Description)
	}

	var res tokenResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return tokenResponse{}, err
	}
	return res, nil
}

func (p *Provider) Exchange(ctx context.Context, code, codeVerifier string) (provider.Token, provider.AccountInfo, error) {
	res, err := p.postTokenForm(ctx, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {p.cfg.RedirectURI},
		"code_verifier": {codeVerifier},
	})
	if err != nil {
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
	res, err := p.postTokenForm(ctx, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token.RefreshToken},
	})
	if err != nil {
		return provider.Token{}, err
	}

	return provider.Token{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(res.ExpiresIn) * time.Second),
		Scope:        res.Scope,
	}, nil
}
