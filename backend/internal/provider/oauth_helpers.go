package provider

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateState returns a random CSRF state value for an OAuth redirect.
func GenerateState() (string, error) {
	return randomURLSafe(32)
}

// GenerateCodeVerifier returns a random PKCE code_verifier.
func GenerateCodeVerifier() (string, error) {
	return randomURLSafe(48)
}

// CodeChallengeS256 derives the PKCE code_challenge (method S256) for a
// given code_verifier, per RFC 7636.
func CodeChallengeS256(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func randomURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
