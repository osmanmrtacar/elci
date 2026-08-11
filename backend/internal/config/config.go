package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	ServerAddr   string
	FrontendURL  string
	DatabasePath string

	JWTSecret string
	JWTTTL    time.Duration

	CORSAllowedOrigins []string

	TikTok    TikTokConfig
	X         XConfig
	Instagram InstagramConfig
	R2        R2Config
}

type TikTokConfig struct {
	ClientKey    string
	ClientSecret string
	RedirectURI  string
}

func (c TikTokConfig) Configured() bool {
	return c.ClientKey != "" && c.ClientSecret != "" && c.RedirectURI != ""
}

type XConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func (c XConfig) Configured() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RedirectURI != ""
}

type InstagramConfig struct {
	AppID       string
	AppSecret   string
	RedirectURI string
}

func (c InstagramConfig) Configured() bool {
	return c.AppID != "" && c.AppSecret != "" && c.RedirectURI != ""
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	PublicBaseURL   string
}

func (c R2Config) Configured() bool {
	return c.AccountID != "" && c.AccessKeyID != "" && c.SecretAccessKey != "" && c.Bucket != ""
}

func Load() (Config, error) {
	cfg := Config{
		ServerAddr:   getenv("SERVER_ADDR", ":8080"),
		FrontendURL:  getenv("FRONTEND_URL", "http://localhost:3000"),
		DatabasePath: getenv("DATABASE_PATH", "./data/elci.db"),

		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTTTL:    30 * 24 * time.Hour,

		CORSAllowedOrigins: splitCSV(getenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),

		TikTok: TikTokConfig{
			ClientKey:    os.Getenv("TIKTOK_CLIENT_KEY"),
			ClientSecret: os.Getenv("TIKTOK_CLIENT_SECRET"),
			RedirectURI:  os.Getenv("TIKTOK_REDIRECT_URI"),
		},
		X: XConfig{
			ClientID:     os.Getenv("X_CLIENT_ID"),
			ClientSecret: os.Getenv("X_CLIENT_SECRET"),
			RedirectURI:  os.Getenv("X_REDIRECT_URI"),
		},
		Instagram: InstagramConfig{
			AppID:       os.Getenv("INSTAGRAM_APP_ID"),
			AppSecret:   os.Getenv("INSTAGRAM_APP_SECRET"),
			RedirectURI: os.Getenv("INSTAGRAM_REDIRECT_URI"),
		},
		R2: R2Config{
			AccountID:       os.Getenv("R2_ACCOUNT_ID"),
			AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
			Bucket:          os.Getenv("R2_BUCKET"),
			PublicBaseURL:   os.Getenv("R2_PUBLIC_BASE_URL"),
		},
	}

	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must be set and at least 32 characters")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
