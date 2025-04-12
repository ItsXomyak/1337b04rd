package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/services"
	"1337b04rd/internal/domain/errors"
)

type ThreadHandler struct {
	threadSvc *services.ThreadService
}

func NewThreadHandler(threadSvc *services.ThreadService) *ThreadHandler {
	return &ThreadHandler{
		threadSvc: threadSvc,
	}
}

func (h *ThreadHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error("failed to parse form", "error", err)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	sessionIDStr := r.FormValue("session_id")
	imageURL := r.FormValue("image_url")
	var imageURLPtr *string
	if imageURL != "" {
		imageURLPtr = &imageURL
	}

	sessionID, err := utils.ParseUUID(sessionIDStr)
	if err != nil {
		logger.Error("invalid session_id", "error", err)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid session ID"})
		return
	}

	if title == "" || content == "" {
		logger.Warn("missing required fields", "title", title, "content", content)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Title and content are required"})
		return
	}

	thread, err := h.threadSvc.CreateThread(ctx, title, content, imageURLPtr, sessionID)
	if err != nil {
		logger.Error("failed to create thread", "error", err)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create thread"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(thread)
}

func (h *ThreadHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		Respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	ctx := r.Context()
	path := r.URL.Path
	if !strings.HasPrefix(path, "/threads/") {
		logger.Warn("invalid path", "path", path)
		Respond(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}

	idStr := strings.TrimPrefix(path, "/threads/")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		logger.Error("invalid thread id", "error", err, "id", idStr)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "Invalid thread ID"})
		return
	}

	thread, err := h.threadSvc.GetThreadByID(ctx, id)
	if err != nil {
		if err == errors.ErrThreadNotFound {
			logger.Warn("thread not found", "id", id)
			Respond(w, http.StatusNotFound, map[string]string{"error": "Thread not found"})
			return
		}
		logger.Error("failed to get thread", "error", err, "id", id)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get thread"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(thread); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

func (h *ThreadHandler) ListActiveThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		Respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	if r.URL.Path != "/threads" {
		logger.Warn("invalid path", "path", r.URL.Path)
		Respond(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}

	ctx := r.Context()
	threads, err := h.threadSvc.ListActiveThreads(ctx)
	if err != nil {
		logger.Error("failed to list active threads", "error", err)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list threads"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(threads); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

func (h *ThreadHandler) ListAllThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		Respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	if r.URL.Path != "/threads/all" {
		logger.Warn("invalid path", "path", r.URL.Path)
		Respond(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}

	ctx := r.Context()
	threads, err := h.threadSvc.ListAllThreads(ctx)
	if err != nil {
		logger.Error("failed to list all threads", "error", err)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list threads"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(threads); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}
