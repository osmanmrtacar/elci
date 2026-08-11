package service

import (
	"context"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/domain"
	"github.com/osmanmertacar/elci/backend/internal/provider"
	"github.com/osmanmertacar/elci/backend/internal/repository"
)

// refreshTokenIfNeeded refreshes a connection's token when it's within 2
// minutes of expiring (or has no known expiry), persisting the result so
// the next call doesn't repeat the refresh. Shared by PostService (before
// publishing) and ConnectionService (before a live AccountInfo lookup).
func refreshTokenIfNeeded(ctx context.Context, p provider.Provider, connections repository.ConnectionRepository, conn domain.PlatformConnection) (provider.Token, error) {
	token := connToken(conn)
	if token.ExpiresAt.IsZero() || time.Until(token.ExpiresAt) > 2*time.Minute {
		return token, nil
	}

	refreshed, err := p.RefreshToken(ctx, token)
	if err != nil {
		return provider.Token{}, err
	}
	conn.AccessToken = refreshed.AccessToken
	conn.RefreshToken = refreshed.RefreshToken
	conn.TokenExpiresAt = refreshed.ExpiresAt
	conn.Scope = refreshed.Scope
	if _, err := connections.Upsert(ctx, conn); err != nil {
		return provider.Token{}, err
	}
	return refreshed, nil
}

func connToken(c domain.PlatformConnection) provider.Token {
	return provider.Token{
		AccessToken:  c.AccessToken,
		RefreshToken: c.RefreshToken,
		ExpiresAt:    c.TokenExpiresAt,
		Scope:        c.Scope,
	}
}
