// Package tiktok implements provider.Provider against TikTok's Content
// Posting API, following TikTok's Direct Post developer guidelines: creator
// info is queried fresh before every post, privacy/interaction settings are
// never defaulted, and media is published via PULL_FROM_URL from our own
// (ownership-verified) storage domain rather than re-uploading files
// ourselves.
package tiktok

import (
	"net/url"
)

type Config struct {
	ClientKey    string
	ClientSecret string
	RedirectURI  string
}

type Provider struct {
	cfg Config
}

func New(cfg Config) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Identifier() string { return "tiktok" }

func (p *Provider) UsesPKCE() bool { return false }

func (p *Provider) AuthURL(state, _ string) string {
	q := url.Values{
		"client_key":    {p.cfg.ClientKey},
		"scope":         {"user.info.basic,video.publish"},
		"response_type": {"code"},
		"redirect_uri":  {p.cfg.RedirectURI},
		"state":         {state},
	}
	return "https://www.tiktok.com/v2/auth/authorize/?" + q.Encode()
}
