package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/osmanmertacar/elci/backend/internal/api/middleware"
	"github.com/osmanmertacar/elci/backend/internal/service"
)

type OAuthHandler struct {
	connections *service.ConnectionService
	frontendURL string
	logger      *slog.Logger
}

func NewOAuthHandler(connections *service.ConnectionService, frontendURL string, logger *slog.Logger) *OAuthHandler {
	return &OAuthHandler{connections: connections, frontendURL: strings.TrimRight(frontendURL, "/"), logger: logger}
}

func (h *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	platform := r.PathValue("platform")

	var userID *int64
	if uid, ok := middleware.UserID(r.Context()); ok {
		userID = &uid
	}

	redirectURL, err := h.connections.Begin(r.Context(), platform, userID)
	if err != nil {
		if errors.Is(err, service.ErrUnsupportedPlatform) {
			writeError(w, http.StatusNotFound, "unsupported platform")
			return
		}
		h.logger.Error("oauth begin failed", "platform", platform, "error", err)
		writeError(w, http.StatusInternalServerError, "could not start sign-in")
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	platform := r.PathValue("platform")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		http.Redirect(w, r, h.frontendURL+"/login?error=missing_code_or_state", http.StatusFound)
		return
	}

	jwtToken, _, err := h.connections.Complete(r.Context(), platform, code, state)
	if err != nil {
		h.logger.Error("oauth callback failed", "platform", platform, "error", err)
		http.Redirect(w, r, h.frontendURL+"/login?error=connection_failed", http.StatusFound)
		return
	}

	http.Redirect(w, r, h.frontendURL+"/callback?token="+url.QueryEscape(jwtToken), http.StatusFound)
}
