package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/api/middleware"
	"github.com/osmanmertacar/elci/backend/internal/domain"
	"github.com/osmanmertacar/elci/backend/internal/repository"
	"github.com/osmanmertacar/elci/backend/internal/service"
)

type MeHandler struct {
	connections repository.ConnectionRepository
	connSvc     *service.ConnectionService
}

func NewMeHandler(connections repository.ConnectionRepository, connSvc *service.ConnectionService) *MeHandler {
	return &MeHandler{connections: connections, connSvc: connSvc}
}

type connectionDTO struct {
	Platform    string    `json:"platform"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	AvatarURL   string    `json:"avatarUrl"`
	IsActive    bool      `json:"isActive"`
	ConnectedAt time.Time `json:"connectedAt"`
}

func (h *MeHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	conns, err := h.connections.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load account")
		return
	}

	dtos := make([]connectionDTO, 0, len(conns))
	for _, c := range conns {
		dtos = append(dtos, connectionDTO{
			Platform:    string(c.Platform),
			Username:    c.Username,
			DisplayName: c.DisplayName,
			AvatarURL:   c.AvatarURL,
			IsActive:    c.IsActive,
			ConnectedAt: c.ConnectedAt,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"userId":      userID,
		"connections": dtos,
	})
}

func (h *MeHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	platform := r.PathValue("platform")
	if err := h.connections.Deactivate(r.Context(), userID, domain.Platform(platform)); err != nil {
		writeError(w, http.StatusInternalServerError, "could not disconnect platform")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ConnectionInfo returns a platform's *live* account data — the composer
// calls this when a platform is selected so it can build TikTok's privacy/
// duet/stitch controls from the account's current capabilities.
func (h *MeHandler) ConnectionInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	platform := r.PathValue("platform")
	info, err := h.connSvc.AccountInfo(r.Context(), userID, platform)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not connected")
		return
	}
	if errors.Is(err, service.ErrUnsupportedPlatform) {
		writeError(w, http.StatusNotFound, "unsupported platform")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadGateway, "could not load live account info: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"platformUserId": info.PlatformUserID,
		"username":       info.Username,
		"displayName":    info.DisplayName,
		"avatarUrl":      info.AvatarURL,
		"extra":          info.Extra,
	})
}
