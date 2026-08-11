package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/osmanmertacar/elci/backend/internal/api/middleware"
	"github.com/osmanmertacar/elci/backend/internal/domain"
	"github.com/osmanmertacar/elci/backend/internal/repository"
	"github.com/osmanmertacar/elci/backend/internal/service"
)

type PostHandler struct {
	posts   repository.PostRepository
	targets repository.PostTargetRepository
	svc     *service.PostService
}

func NewPostHandler(posts repository.PostRepository, targets repository.PostTargetRepository, svc *service.PostService) *PostHandler {
	return &PostHandler{posts: posts, targets: targets, svc: svc}
}

type targetInputRequest struct {
	Platform          string         `json:"platform"`
	CaptionOverride   *string        `json:"captionOverride"`
	MediaKindOverride *string        `json:"mediaKindOverride"`
	MediaURLsOverride []string       `json:"mediaUrlsOverride"`
	Settings          map[string]any `json:"settings"`
}

type createPostRequest struct {
	DefaultCaption   string               `json:"defaultCaption"`
	DefaultMediaKind *string              `json:"defaultMediaKind"`
	DefaultMediaURLs []string             `json:"defaultMediaUrls"`
	Targets          []targetInputRequest `json:"targets"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := service.CreatePostInput{
		DefaultCaption:   req.DefaultCaption,
		DefaultMediaKind: parseMediaKind(req.DefaultMediaKind),
		DefaultMediaURLs: req.DefaultMediaURLs,
	}
	for _, t := range req.Targets {
		input.Targets = append(input.Targets, service.TargetInput{
			Platform:          domain.Platform(t.Platform),
			CaptionOverride:   t.CaptionOverride,
			MediaKindOverride: parseMediaKind(t.MediaKindOverride),
			MediaURLsOverride: t.MediaURLsOverride,
			Settings:          t.Settings,
		})
	}

	post, targets, err := h.svc.CreateDraft(r.Context(), userID, input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, postResponse(post, targets))
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	posts, err := h.posts.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list posts")
		return
	}

	out := make([]map[string]any, 0, len(posts))
	for _, p := range posts {
		targets, err := h.targets.ListByPost(r.Context(), p.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not load post targets")
			return
		}
		out = append(out, postResponse(p, targets))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.posts.Get(r.Context(), userID, id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load post")
		return
	}
	targets, err := h.targets.ListByPost(r.Context(), post.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load post targets")
		return
	}
	writeJSON(w, http.StatusOK, postResponse(post, targets))
}

type scheduleRequest struct {
	ScheduledAt time.Time `json:"scheduledAt"`
}

func (h *PostHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	var req scheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ScheduledAt.IsZero() {
		writeError(w, http.StatusBadRequest, "scheduledAt is required")
		return
	}

	if err := h.svc.Schedule(r.Context(), userID, id, req.ScheduledAt); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) Publish(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	if err := h.svc.PublishNow(r.Context(), userID, id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}
	if err := h.posts.Delete(r.Context(), userID, id); err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete post")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error) {
	var vf *service.ValidationFailure
	if errors.As(err, &vf) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": vf.ByTarget})
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func parseMediaKind(s *string) *domain.MediaKind {
	if s == nil {
		return nil
	}
	k := domain.MediaKind(*s)
	return &k
}

func postResponse(p domain.Post, targets []domain.PostTarget) map[string]any {
	targetsOut := make([]map[string]any, 0, len(targets))
	for _, t := range targets {
		targetsOut = append(targetsOut, map[string]any{
			"id":                t.ID,
			"platform":          t.Platform,
			"captionOverride":   t.CaptionOverride,
			"mediaUrlsOverride": t.MediaURLsOverride,
			"settings":          t.Settings,
			"status":            t.Status,
			"platformPostId":    t.PlatformPostID,
			"error":             t.Error,
			"publishedAt":       t.PublishedAt,
		})
	}
	return map[string]any{
		"id":               p.ID,
		"defaultCaption":   p.DefaultCaption,
		"defaultMediaKind": p.DefaultMediaKind,
		"defaultMediaUrls": p.DefaultMediaURLs,
		"status":           p.Status,
		"scheduledAt":      p.ScheduledAt,
		"createdAt":        p.CreatedAt,
		"targets":          targetsOut,
	}
}
