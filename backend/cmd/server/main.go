package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/osmanmertacar/elci/backend/internal/api"
	"github.com/osmanmertacar/elci/backend/internal/auth"
	"github.com/osmanmertacar/elci/backend/internal/config"
	"github.com/osmanmertacar/elci/backend/internal/database"
	"github.com/osmanmertacar/elci/backend/internal/provider"
	"github.com/osmanmertacar/elci/backend/internal/provider/instagram"
	"github.com/osmanmertacar/elci/backend/internal/provider/tiktok"
	"github.com/osmanmertacar/elci/backend/internal/provider/x"
	"github.com/osmanmertacar/elci/backend/internal/repository"
	"github.com/osmanmertacar/elci/backend/internal/service"
	"github.com/osmanmertacar/elci/backend/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// In Docker, env vars are injected directly and no .env file exists —
	// that's fine, godotenv only supplements what's already in the process env.
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid config", "error", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		logger.Error("could not open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	users := repository.NewUserRepository(db)
	connections := repository.NewConnectionRepository(db)
	sessions := repository.NewOAuthSessionRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	postRepo := repository.NewPostRepository(db)
	postTargetRepo := repository.NewPostTargetRepository(db)

	var mediaService *service.MediaService
	if cfg.R2.Configured() {
		mediaService, err = service.NewMediaService(context.Background(), cfg.R2)
		if err != nil {
			logger.Error("could not init media service", "error", err)
			os.Exit(1)
		}
	} else {
		logger.Warn("R2 not configured, media upload disabled")
	}

	registry := provider.NewRegistry()
	if cfg.TikTok.Configured() {
		registry.Register(tiktok.New(tiktok.Config{
			ClientKey:    cfg.TikTok.ClientKey,
			ClientSecret: cfg.TikTok.ClientSecret,
			RedirectURI:  cfg.TikTok.RedirectURI,
		}))
	} else {
		logger.Warn("TikTok not configured, skipping registration")
	}
	if cfg.X.Configured() {
		registry.Register(x.New(x.Config{
			ClientID:     cfg.X.ClientID,
			ClientSecret: cfg.X.ClientSecret,
			RedirectURI:  cfg.X.RedirectURI,
		}))
	} else {
		logger.Warn("X not configured, skipping registration")
	}
	if cfg.Instagram.Configured() {
		registry.Register(instagram.New(instagram.Config{
			AppID:       cfg.Instagram.AppID,
			AppSecret:   cfg.Instagram.AppSecret,
			RedirectURI: cfg.Instagram.RedirectURI,
		}))
	} else {
		logger.Warn("Instagram not configured, skipping registration")
	}

	tokens := auth.NewTokenIssuer(cfg.JWTSecret, cfg.JWTTTL)
	connectionService := service.NewConnectionService(registry, users, connections, sessions, tokens)
	postService := service.NewPostService(registry, postRepo, postTargetRepo, connections, mediaRepo)

	handler := api.NewRouter(api.Deps{
		Tokens:             tokens,
		Connections:        connectionService,
		ConnectionRepo:     connections,
		Media:              mediaService,
		MediaRepo:          mediaRepo,
		Posts:              postService,
		PostRepo:           postRepo,
		PostTargetRepo:     postTargetRepo,
		FrontendURL:        cfg.FrontendURL,
		CORSAllowedOrigins: cfg.CORSAllowedOrigins,
		Logger:             logger,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	scheduler := worker.NewScheduler(postRepo, postService, logger)
	statusPoller := worker.NewStatusPoller(postTargetRepo, postService, logger)
	go scheduler.Run(ctx)
	go statusPoller.Run(ctx)

	srv := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server starting", "addr", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
