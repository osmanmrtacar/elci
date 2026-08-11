// Package x implements provider.Provider against X's API v2: OAuth2 with
// PKCE, chunked media upload, and the /2/tweets endpoint. Unlike TikTok, a
// successful POST to /2/tweets publishes immediately — there is no async
// processing step to poll.
package x

import "net/url"

const authorizeURL = "https://x.com/i/oauth2/authorize"

// scopes needed to read the account's own profile and post text+media.
const scopes = "tweet.read tweet.write users.read offline.access media.write"

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type Provider struct {
	cfg Config
}

func New(cfg Config) *Provider {
	return &Provider{cfg: cfg}
}

func (p *Provider) Identifier() string { return "x" }

func (p *Provider) UsesPKCE() bool { return true }

func (p *Provider) AuthURL(state, codeChallenge string) string {
	q := url.Values{
		"response_type":         {"code"},
		"client_id":             {p.cfg.ClientID},
		"redirect_uri":          {p.cfg.RedirectURI},
		"scope":                 {scopes},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}
	return authorizeURL + "?" + q.Encode()
}
