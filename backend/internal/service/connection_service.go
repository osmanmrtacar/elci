package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/auth"
	"github.com/osmanmertacar/elci/backend/internal/domain"
	"github.com/osmanmertacar/elci/backend/internal/provider"
	"github.com/osmanmertacar/elci/backend/internal/repository"
)

var ErrUnsupportedPlatform = errors.New("unsupported platform")

type ConnectionService struct {
	registry    *provider.Registry
	users       repository.UserRepository
	connections repository.ConnectionRepository
	sessions    repository.OAuthSessionRepository
	tokens      auth.TokenIssuer
}

func NewConnectionService(
	registry *provider.Registry,
	users repository.UserRepository,
	connections repository.ConnectionRepository,
	sessions repository.OAuthSessionRepository,
	tokens auth.TokenIssuer,
) *ConnectionService {
	return &ConnectionService{
		registry:    registry,
		users:       users,
		connections: connections,
		sessions:    sessions,
		tokens:      tokens,
	}
}

// Begin starts an OAuth authorization for platform and returns the URL to
// redirect the browser to. userID is non-nil when an already logged-in user
// is connecting an additional platform.
func (s *ConnectionService) Begin(ctx context.Context, platform string, userID *int64) (string, error) {
	p, ok := s.registry.Get(platform)
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedPlatform, platform)
	}

	state, err := provider.GenerateState()
	if err != nil {
		return "", err
	}

	var codeVerifier string
	if p.UsesPKCE() {
		codeVerifier, err = provider.GenerateCodeVerifier()
		if err != nil {
			return "", err
		}
	}

	_, err = s.sessions.Create(ctx, domain.OAuthSession{
		State:        state,
		CodeVerifier: codeVerifier,
		Platform:     domain.Platform(platform),
		UserID:       userID,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	challenge := codeVerifier
	if p.UsesPKCE() {
		challenge = provider.CodeChallengeS256(codeVerifier)
	}

	return p.AuthURL(state, challenge), nil
}

// Complete redeems the OAuth callback, upserts the platform connection, and
// returns a session JWT for the resulting user.
func (s *ConnectionService) Complete(ctx context.Context, platform, code, state string) (jwtToken string, userID int64, err error) {
	session, err := s.sessions.Consume(ctx, state)
	if err != nil {
		return "", 0, err
	}
	if string(session.Platform) != platform {
		return "", 0, errors.New("state does not match platform")
	}

	p, ok := s.registry.Get(platform)
	if !ok {
		return "", 0, fmt.Errorf("%w: %s", ErrUnsupportedPlatform, platform)
	}

	tok, account, err := p.Exchange(ctx, code, session.CodeVerifier)
	if err != nil {
		return "", 0, err
	}

	if session.UserID != nil {
		userID = *session.UserID
	} else {
		existing, err := s.connections.FindByPlatformUserID(ctx, domain.Platform(platform), account.PlatformUserID)
		switch {
		case err == nil:
			userID = existing.UserID
		case errors.Is(err, repository.ErrNotFound):
			newUser, err := s.users.Create(ctx)
			if err != nil {
				return "", 0, err
			}
			userID = newUser.ID
		default:
			return "", 0, err
		}
	}

	_, err = s.connections.Upsert(ctx, domain.PlatformConnection{
		UserID:         userID,
		Platform:       domain.Platform(platform),
		PlatformUserID: account.PlatformUserID,
		Username:       account.Username,
		DisplayName:    account.DisplayName,
		AvatarURL:      account.AvatarURL,
		AccessToken:    tok.AccessToken,
		RefreshToken:   tok.RefreshToken,
		TokenExpiresAt: tok.ExpiresAt,
		Scope:          tok.Scope,
	})
	if err != nil {
		return "", 0, err
	}

	jwtToken, err = s.tokens.Issue(userID)
	if err != nil {
		return "", 0, err
	}

	return jwtToken, userID, nil
}

// AccountInfo fetches a connection's *live* platform data (e.g. TikTok's
// creator_info: privacy options, duet/stitch/comment availability) — the
// composer calls this on selecting a platform so it never shows stale
// capabilities from whenever the account was first connected.
func (s *ConnectionService) AccountInfo(ctx context.Context, userID int64, platform string) (provider.AccountInfo, error) {
	p, ok := s.registry.Get(platform)
	if !ok {
		return provider.AccountInfo{}, fmt.Errorf("%w: %s", ErrUnsupportedPlatform, platform)
	}

	conn, err := s.connections.Get(ctx, userID, domain.Platform(platform))
	if err != nil {
		return provider.AccountInfo{}, err
	}

	token, err := refreshTokenIfNeeded(ctx, p, s.connections, conn)
	if err != nil {
		return provider.AccountInfo{}, err
	}

	return p.AccountInfo(ctx, token)
}
