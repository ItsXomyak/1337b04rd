package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/services"
	"1337b04rd/internal/domain/errors"
)

type CommentHandler struct {
	commentSvc *services.CommentService
}

func NewCommentHandler(commentSvc *services.CommentService, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{
		commentSvc: commentSvc,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		Respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	ctx := r.Context()
	path := r.URL.Path
	if !strings.HasPrefix(path, "/threads/") || !strings.Contains(path, "/comments") {
		logger.Warn("invalid path", "path", path)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}

	//	/threads/{thread_id}/comments
	parts := strings.Split(strings.TrimPrefix(path, "/threads/"), "/")
	if len(parts) < 2 || parts[1] != "comments" {
		logger.Warn("invalid path format", "path", path)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	threadIDStr := parts[0]
	threadID, err := utils.ParseUUID(threadIDStr)
	if err != nil {
		logger.Error("invalid thread_id", "error", err, "thread_id", threadIDStr)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid thread ID"})
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		logger.Error("failed to parse form", "error", err)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
		return
	}

	content := r.FormValue("content")
	sessionIDStr := r.FormValue("session_id")
	parentIDStr := r.FormValue("parent_id")
	imageURL := r.FormValue("image_url")

	var parentID *utils.UUID
	if parentIDStr != "" {
		parsedID, err := utils.ParseUUID(parentIDStr)
		if err != nil {
			logger.Error("invalid parent_id", "error", err, "parent_id", parentIDStr)
			Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid parent_id"})
			return
		}
		parentID = &parsedID
	}

	var imageURLPtr *string
	if imageURL != "" {
		imageURLPtr = &imageURL
	}

	sessionID, err := utils.ParseUUID(sessionIDStr)
	if err != nil {
		logger.Error("invalid session_id", "error", err)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid session_id"})
		return
	}

	if content == "" {
		logger.Warn("missing content", "thread_id", threadID)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Content is required"})
		return
	}

	comment, err := h.commentSvc.CreateComment(ctx, threadID, parentID, content, imageURLPtr, sessionID)
	if err != nil {
		if err == errors.ErrThreadNotFound {
			logger.Warn("thread not found", "thread_id", threadID)
			Respond(w, http.StatusNotFound, map[string]string{"error": "Thread not found"})
			return
		}
		logger.Error("failed to create comment", "error", err, "thread_id", threadID)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create comment"})
		return
	}

	Respond(w, http.StatusCreated, comment)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "Failed to encode response"})
	}
}

func (h *CommentHandler) GetCommentsByThreadID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		Respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	ctx := r.Context()
	path := r.URL.Path
	if !strings.HasPrefix(path, "/threads/") || !strings.Contains(path, "/comments") {
		logger.Warn("invalid path", "path", path)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}

	// Извлекаем thread_id из пути: /threads/{thread_id}/comments
	parts := strings.Split(strings.TrimPrefix(path, "/threads/"), "/")
	if len(parts) < 2 || parts[1] != "comments" {
		logger.Warn("invalid path format", "path", path)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
		return
	}
	threadIDStr := parts[0]
	threadID, err := utils.ParseUUID(threadIDStr)
	if err != nil {
		logger.Error("invalid thread_id", "error", err, "thread_id", threadIDStr)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid thread ID"})
		return
	}

	comments, err := h.commentSvc.GetCommentsByThreadID(ctx, threadID)
	if err != nil {
		if err == errors.ErrThreadNotFound {
			logger.Warn("thread not found", "thread_id", threadID)
			Respond(w, http.StatusNotFound, map[string]string{"error": "Thread not found"})
			return
		}
		logger.Error("failed to get comments", "error", err, "thread_id", threadID)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get comments"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}
