package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	TikTok    TikTokConfig
	X         XConfig
	Instagram InstagramConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	Log       LogConfig
}

type ServerConfig struct {
	Port        string
	Host        string
	FrontendURL string
	Environment string
}

type TikTokConfig struct {
	ClientKey    string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

type XConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type InstagramConfig struct {
	AppID       string
	AppSecret   string
	RedirectURI string
}

type DatabaseConfig struct {
	Path string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

type LogConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", "localhost"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		TikTok: TikTokConfig{
			ClientKey:    getEnv("TIKTOK_CLIENT_KEY", ""),
			ClientSecret: getEnv("TIKTOK_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("TIKTOK_REDIRECT_URI", ""),
			Scopes:       strings.Split(getEnv("TIKTOK_SCOPES", "user.info.basic,video.publish"), ","),
		},
		X: XConfig{
			ClientID:     getEnv("X_CLIENT_ID", ""),
			ClientSecret: getEnv("X_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("X_REDIRECT_URI", ""),
		},
		Instagram: InstagramConfig{
			AppID:       getEnv("INSTAGRAM_APP_ID", ""),
			AppSecret:   getEnv("INSTAGRAM_APP_SECRET", ""),
			RedirectURI: getEnv("INSTAGRAM_REDIRECT_URI", ""),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/sosyal.db"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Expiration: parseDuration(getEnv("JWT_EXPIRATION", "24h")),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if required configuration fields are set
func (c *Config) Validate() error {
	// Check if at least one platform is configured
	hasTikTok := c.TikTok.ClientKey != "" && c.TikTok.ClientSecret != "" && c.TikTok.RedirectURI != ""
	hasX := c.X.ClientID != "" && c.X.ClientSecret != "" && c.X.RedirectURI != ""
	hasInstagram := c.Instagram.AppID != "" && c.Instagram.AppSecret != "" && c.Instagram.RedirectURI != ""

	if !hasTikTok && !hasX && !hasInstagram {
		return fmt.Errorf("at least one platform (TikTok, X, or Instagram) must be fully configured")
	}

	// Validate TikTok config if any TikTok field is set
	if c.TikTok.ClientKey != "" || c.TikTok.ClientSecret != "" || c.TikTok.RedirectURI != "" {
		if c.TikTok.ClientKey == "" {
			return fmt.Errorf("TIKTOK_CLIENT_KEY is required when TikTok is configured")
		}
		if c.TikTok.ClientSecret == "" {
			return fmt.Errorf("TIKTOK_CLIENT_SECRET is required when TikTok is configured")
		}
		if c.TikTok.RedirectURI == "" {
			return fmt.Errorf("TIKTOK_REDIRECT_URI is required when TikTok is configured")
		}
	}

	// Validate X config if any X field is set
	if c.X.ClientID != "" || c.X.ClientSecret != "" || c.X.RedirectURI != "" {
		if c.X.ClientID == "" {
			return fmt.Errorf("X_CLIENT_ID is required when X is configured")
		}
		if c.X.ClientSecret == "" {
			return fmt.Errorf("X_CLIENT_SECRET is required when X is configured")
		}
		if c.X.RedirectURI == "" {
			return fmt.Errorf("X_REDIRECT_URI is required when X is configured")
		}
	}

	// Validate Instagram config if any Instagram field is set
	if c.Instagram.AppID != "" || c.Instagram.AppSecret != "" || c.Instagram.RedirectURI != "" {
		if c.Instagram.AppID == "" {
			return fmt.Errorf("INSTAGRAM_APP_ID is required when Instagram is configured")
		}
		if c.Instagram.AppSecret == "" {
			return fmt.Errorf("INSTAGRAM_APP_SECRET is required when Instagram is configured")
		}
		if c.Instagram.RedirectURI == "" {
			return fmt.Errorf("INSTAGRAM_REDIRECT_URI is required when Instagram is configured")
		}
	}

	// JWT is always required
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// parseDuration parses a duration string, returns 24h as default on error
func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return duration
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
