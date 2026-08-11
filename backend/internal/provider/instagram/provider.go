package instagram

import "net/url"

const authorizeURL = "https://www.instagram.com/oauth/authorize"

const scopes = "instagram_business_basic,instagram_business_content_publish"

type Config struct {
	AppID       string
	AppSecret   string
	RedirectURI string
}

type Provider struct {
	cfg Config
}

func New(cfg Config) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Identifier() string { return "instagram" }

func (p *Provider) UsesPKCE() bool { return false }

func (p *Provider) AuthURL(state, _ string) string {
	q := url.Values{
		"client_id":     {p.cfg.AppID},
		"redirect_uri":  {p.cfg.RedirectURI},
		"response_type": {"code"},
		"scope":         {scopes},
		"state":         {state},
	}
	return authorizeURL + "?" + q.Encode()
}
