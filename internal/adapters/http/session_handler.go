package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/services"
)

type SessionHandler struct {
	SessionService *services.SessionService
}

type changeNameRequest struct {
	DisplayName string `json:"display_name"`
}

type changeNameResponse struct {
	Success bool `json:"success"`
}

type sessionInfoResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type sessionItem struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	ExpiresAt   string `json:"expires_at"`
}

// POST /session/name
func (h *SessionHandler) ChangeDisplayName(w http.ResponseWriter, r *http.Request) {
	sess, ok := GetSessionFromContext(r.Context())
	if !ok {
		logger.Warn("session not found in context")
		Respond(w, http.StatusUnauthorized, map[string]string{"error": "session not found"})
		return
	}

	var req changeNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("failed to decode change name request", "err", err)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	name := strings.TrimSpace(req.DisplayName)
	if len(name) < 2 || len(name) > 30 {
		logger.Warn("invalid name length", "name", name)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "invalid name length"})
		return
	}

	err := h.SessionService.UpdateDisplayName(sess.ID, name)
	if err != nil {
		logger.Error("failed to update display name", "session_id", sess.ID, "err", err)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "could not update name"})
		return
	}

	logger.Info("display name updated", "session_id", sess.ID, "new_name", name)
	Respond(w, http.StatusOK, changeNameResponse{Success: true})
}

// GET /session/me
func (h *SessionHandler) GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	sess, ok := GetSessionFromContext(r.Context())
	if !ok {
		logger.Warn("session not found in context (me endpoint)")
		Respond(w, http.StatusUnauthorized, map[string]string{"error": "session not found"})
		return
	}

	logger.Info("session info fetched", "session_id", sess.ID)
	resp := sessionInfoResponse{
		ID:          sess.ID.String(),
		DisplayName: sess.DisplayName,
		AvatarURL:   sess.AvatarURL,
	}

	Respond(w, http.StatusOK, resp)
}

// GET /session/list
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.SessionService.ListActiveSessions()
	if err != nil {
		logger.Error("failed to list active sessions", "err", err)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "could not list sessions"})
		return
	}

	logger.Info("active sessions listed", "count", len(sessions))

	result := make([]sessionItem, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, sessionItem{
			ID:          s.ID.String(),
			DisplayName: s.DisplayName,
			AvatarURL:   s.AvatarURL,
			ExpiresAt:   s.ExpiresAt.Format("2006-01-02 15:04:05"),
		})
	}

	Respond(w, http.StatusOK, result)
}
