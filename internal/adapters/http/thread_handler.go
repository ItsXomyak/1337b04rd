package http

import (
	"encoding/json"
	"io"
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
	sess, ok := GetSessionFromContext(r.Context())
	if !ok {
		logger.Warn("session not found in CreateThread")
		Respond(w, http.StatusUnauthorized, map[string]string{"error": "session not found"})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Error("failed to parse multipart form", "error", err)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	if title == "" || content == "" {
		Respond(w, http.StatusBadRequest, map[string]string{"error": "title and content are required"})
		return
	}

	var file io.Reader
	var contentType string

	imageHeader, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer imageHeader.Close()
		contentType = fileHeader.Header.Get("Content-Type")
		file = imageHeader
	} else if err != http.ErrMissingFile {
		logger.Error("error retrieving file from form", "error", err)
		Respond(w, http.StatusBadRequest, map[string]string{"error": "invalid image upload"})
		return
	}

	t, err := h.threadSvc.CreateThread(r.Context(), title, content, file, contentType, sess.ID)
	if err != nil {
		logger.Error("failed to create thread", "error", err)
		Respond(w, http.StatusInternalServerError, map[string]string{"error": "could not create thread"})
		return
	}

	Respond(w, http.StatusCreated, map[string]string{
		"thread_id": t.ID.String(),
	})
}

func (h *ThreadHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		Respond(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	ctx := r.Context()
	path := r.URL.Path
	if !strings.HasPrefix(path, "/threads/view/") {
		logger.Warn("invalid path", "path", path)
		Respond(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}

	idStr := strings.TrimPrefix(path, "/threads/view/")
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
