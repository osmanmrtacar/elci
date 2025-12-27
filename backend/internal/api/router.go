package api

import (
	"github.com/gin-gonic/gin"
	"github.com/osmanmertacar/sosyal/backend/internal/api/handlers"
	"github.com/osmanmertacar/sosyal/backend/internal/api/middleware"
	"github.com/osmanmertacar/sosyal/backend/internal/config"
	"github.com/osmanmertacar/sosyal/backend/internal/database"
	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
	"github.com/osmanmertacar/sosyal/backend/internal/services"
	"github.com/osmanmertacar/sosyal/backend/internal/services/platform"
)

// SetupRouter sets up the HTTP router with all routes
func SetupRouter(cfg *config.Config, db *database.DB) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Trust all proxies (removes warning, assumes Cloudflare strips malicious headers)
	// This is safe when using Cloudflare proxy as it sanitizes X-Forwarded-For
	router.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})

	// Apply CORS middleware
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		if err := db.Health(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Initialize repositories
	userRepo := models.NewUserRepository(db.DB)
	tokenRepo := models.NewTokenRepository(db.DB)
	postRepo := models.NewPostRepository(db.DB)
	platformConnectionRepo := models.NewPlatformConnectionRepository(db.DB)
	oauthSessionRepo := models.NewOAuthSessionRepository(db.DB)

	// Initialize platform registry
	platformRegistry := platform.NewPlatformRegistry()

	// Initialize TikTok platform services
	tiktokService := services.NewTikTokService(cfg)
	tokenService := services.NewTokenService(cfg, tokenRepo)
	tokenService.SetTikTokService(tiktokService)
	authService := services.NewAuthService(cfg, tokenService, tiktokService, userRepo)

	// Register TikTok platform
	tiktokPlatform := platform.NewTikTokPlatformService(cfg, tiktokService)
	platformRegistry.Register(tiktokPlatform)

	// Initialize X platform services (if configured)
	if cfg.X.ClientID != "" && cfg.X.ClientSecret != "" {
		xPlatform := platform.NewXPlatformService(
			cfg.X.ClientID,
			cfg.X.ClientSecret,
			cfg.X.RedirectURI,
		)
		platformRegistry.Register(xPlatform)
	}

	// Initialize Instagram platform services (if configured)
	if cfg.Instagram.AppID != "" && cfg.Instagram.AppSecret != "" {
		instagramPlatform := platform.NewInstagramPlatformService(
			cfg.Instagram.AppID,
			cfg.Instagram.AppSecret,
			cfg.Instagram.RedirectURI,
		)
		platformRegistry.Register(instagramPlatform)
	}

	// (postService kept for potential backward compatibility if needed)

	// Initialize multi-platform post service
	multiPlatformPostService := services.NewMultiPlatformPostService(
		postRepo,
		tokenRepo,
		platformConnectionRepo,
		platformRegistry,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userRepo)
	multiPlatformAuthHandler := handlers.NewMultiPlatformAuthHandler(
		cfg,
		platformRegistry,
		userRepo,
		tokenRepo,
		platformConnectionRepo,
		oauthSessionRepo,
	)
	multiPlatformPostHandler := handlers.NewMultiPlatformPostHandler(multiPlatformPostService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no auth middleware)
		auth := v1.Group("/auth")
		{
			// Multi-platform auth routes
			auth.GET("/tiktok/login", multiPlatformAuthHandler.TikTokLogin)
			auth.GET("/tiktok/callback", multiPlatformAuthHandler.TikTokCallback)
			auth.GET("/x/login", multiPlatformAuthHandler.XLogin)
			auth.GET("/x/callback", multiPlatformAuthHandler.XCallback)
			auth.GET("/instagram/login", multiPlatformAuthHandler.InstagramLogin)
			auth.GET("/instagram/callback", multiPlatformAuthHandler.InstagramCallback)
			auth.POST("/logout", multiPlatformAuthHandler.Logout)

			// Development only: Get JWT token for testing
			if cfg.Server.Environment == "development" {
				auth.GET("/dev/token/:user_id", authHandler.GetDevToken)
			}
		}

		// Protected routes (require auth)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// User routes
			protected.GET("/auth/me", multiPlatformAuthHandler.GetCurrentUser)
			protected.GET("/auth/platforms", multiPlatformAuthHandler.GetConnectedPlatforms)
			protected.DELETE("/auth/platforms/:platform", multiPlatformAuthHandler.DisconnectPlatform)

			// Post routes - using multi-platform handler
			posts := protected.Group("/posts")
			{
				posts.POST("", multiPlatformPostHandler.CreatePost)
				posts.GET("", multiPlatformPostHandler.GetPosts)
				posts.GET("/:id", multiPlatformPostHandler.GetPost)
				posts.GET("/:id/status", multiPlatformPostHandler.GetPostStatus)
			}
		}
	}

	return router
}
