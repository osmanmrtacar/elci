package api

import (
	"log/slog"
	"net/http"

	"github.com/osmanmertacar/elci/backend/internal/api/handlers"
	"github.com/osmanmertacar/elci/backend/internal/api/middleware"
	"github.com/osmanmertacar/elci/backend/internal/auth"
	"github.com/osmanmertacar/elci/backend/internal/repository"
	"github.com/osmanmertacar/elci/backend/internal/service"
)

type Deps struct {
	Tokens             auth.TokenIssuer
	Connections        *service.ConnectionService
	ConnectionRepo     repository.ConnectionRepository
	Media              *service.MediaService
	MediaRepo          repository.MediaRepository
	Posts              *service.PostService
	PostRepo           repository.PostRepository
	PostTargetRepo     repository.PostTargetRepository
	FrontendURL        string
	CORSAllowedOrigins []string
	Logger             *slog.Logger
}

func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	oauthHandler := handlers.NewOAuthHandler(deps.Connections, deps.FrontendURL, deps.Logger)
	meHandler := handlers.NewMeHandler(deps.ConnectionRepo, deps.Connections)

	optionalAuth := middleware.OptionalAuth(deps.Tokens)
	requireAuth := middleware.RequireAuth(deps.Tokens)

	mux.Handle("GET /api/v1/auth/{platform}/login", optionalAuth(http.HandlerFunc(oauthHandler.Login)))
	mux.HandleFunc("GET /api/v1/auth/{platform}/callback", oauthHandler.Callback)

	mux.Handle("GET /api/v1/me", requireAuth(http.HandlerFunc(meHandler.Me)))
	mux.Handle("DELETE /api/v1/connections/{platform}", requireAuth(http.HandlerFunc(meHandler.Disconnect)))
	mux.Handle("GET /api/v1/connections/{platform}/info", requireAuth(http.HandlerFunc(meHandler.ConnectionInfo)))

	if deps.Media != nil {
		mediaHandler := handlers.NewMediaHandler(deps.Media, deps.MediaRepo)
		mux.Handle("POST /api/v1/media/presign", requireAuth(http.HandlerFunc(mediaHandler.Presign)))
		mux.Handle("POST /api/v1/media/confirm", requireAuth(http.HandlerFunc(mediaHandler.Confirm)))
	}

	postHandler := handlers.NewPostHandler(deps.PostRepo, deps.PostTargetRepo, deps.Posts)
	mux.Handle("POST /api/v1/posts", requireAuth(http.HandlerFunc(postHandler.Create)))
	mux.Handle("GET /api/v1/posts", requireAuth(http.HandlerFunc(postHandler.List)))
	mux.Handle("GET /api/v1/posts/{id}", requireAuth(http.HandlerFunc(postHandler.Get)))
	mux.Handle("POST /api/v1/posts/{id}/schedule", requireAuth(http.HandlerFunc(postHandler.Schedule)))
	mux.Handle("POST /api/v1/posts/{id}/publish", requireAuth(http.HandlerFunc(postHandler.Publish)))
	mux.Handle("DELETE /api/v1/posts/{id}", requireAuth(http.HandlerFunc(postHandler.Delete)))

	var handler http.Handler = mux
	handler = middleware.CORS(deps.CORSAllowedOrigins)(handler)
	handler = middleware.SecurityHeaders(handler)

	return handler
}
