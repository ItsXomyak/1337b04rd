package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/services"
	"1337b04rd/internal/domain/errors"
)

type ThreadHandler struct {
	threadSvc *services.ThreadService
	s3Client  *s3.S3Client
	sessionSvc *services.SessionService
}

func NewThreadHandler(threadSvc *services.ThreadService) *ThreadHandler {
	return &ThreadHandler{
		threadSvc: threadSvc,
	}
}

func (h *ThreadHandler) SubmitPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Warn("failed to parse form", "error", err)
		RenderTemplate(w, "error.html", map[string]interface{}{
			"Code":    http.StatusBadRequest,
			"Message": "Invalid form data",
		})
		return
	}

	sess, ok := GetSessionFromContext(r.Context())
	if !ok {
		logger.Warn("session not found")
		RenderTemplate(w, "error.html", map[string]interface{}{
			"Code":    http.StatusUnauthorized,
			"Message": "Session required",
		})
		return
	}

	title := r.FormValue("subject")
	content := r.FormValue("comment")
	name := r.FormValue("name")
	file, fileHeader, err := r.FormFile("file")
	var imageURL *string
	if err == nil && file != nil && fileHeader != nil {
		defer file.Close()
		contentType := fileHeader.Header.Get("Content-Type")
		url, err := h.s3Client.UploadPostImage(file, fileHeader.Filename, contentType)
		if err != nil {
			logger.Error("failed to upload image to S3", "error", err)
			RenderTemplate(w, "error.html", map[string]interface{}{
				"Code":    http.StatusInternalServerError,
				"Message": "Failed to upload image",
			})
			return
		}
		imageURL = &url
	}

	if name != "" {
		if err := h.sessionSvc.UpdateDisplayName(sess.ID, name); err != nil {
			logger.Error("failed to update display name", "error", err)
			RenderTemplate(w, "error.html", map[string]interface{}{
				"Code":    http.StatusInternalServerError,
				"Message": "Failed to update name",
			})
			return
		}
	}

	thread, err := h.threadSvc.CreateThread(r.Context(), title, content, imageURL, sess.ID)
	if err != nil {
		logger.Error("failed to create thread", "error", err)
		RenderTemplate(w, "error.html", map[string]interface{}{
			"Code":    http.StatusInternalServerError,
			"Message": "Failed to create thread",
		})
		return
	}

	http.Redirect(w, r, "/post/"+thread.ID.String(), http.StatusSeeOther)
}

func (h *ThreadHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error("failed to parse form", "error", err)
		http.Error(w, `{"error": "Invalid form data"}`, http.StatusBadRequest)
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
		http.Error(w, `{"error": "Invalid session_id"}`, http.StatusBadRequest)
		return
	}

	if title == "" || content == "" {
		logger.Warn("missing required fields", "title", title, "content", content)
		http.Error(w, `{"error": "Title and content are required"}`, http.StatusBadRequest)
		return
	}

	thread, err := h.threadSvc.CreateThread(ctx, title, content, imageURLPtr, sessionID)
	if err != nil {
		logger.Error("failed to create thread", "error", err)
		http.Error(w, `{"error": "Failed to create thread"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(thread)
}

func (h *ThreadHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	path := r.URL.Path
	if !strings.HasPrefix(path, "/threads/") {
		logger.Warn("invalid path", "path", path)
		http.Error(w, `{"error": "Invalid thread ID"}`, http.StatusBadRequest)
		return
	}

	idStr := strings.TrimPrefix(path, "/threads/")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		logger.Error("invalid thread id", "error", err, "id", idStr)
		http.Error(w, `{"error": "Invalid thread ID"}`, http.StatusBadRequest)
		return
	}

	thread, err := h.threadSvc.GetThreadByID(ctx, id)
	if err != nil {
		if err == errors.ErrThreadNotFound {
			logger.Warn("thread not found", "id", id)
			http.Error(w, `{"error": "Thread not found"}`, http.StatusNotFound)
			return
		}
		logger.Error("failed to get thread", "error", err, "id", id)
		http.Error(w, `{"error": "Failed to get thread"}`, http.StatusInternalServerError)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/threads" {
		logger.Warn("invalid path", "path", r.URL.Path)
		http.Error(w, `{"error": "Not found"}`, http.StatusNotFound)
		return
	}

	ctx := r.Context()
	threads, err := h.threadSvc.ListActiveThreads(ctx)
	if err != nil {
		logger.Error("failed to list active threads", "error", err)
		http.Error(w, `{"error": "Failed to list threads"}`, http.StatusInternalServerError)
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/threads/all" {
		logger.Warn("invalid path", "path", r.URL.Path)
		http.Error(w, `{"error": "Not found"}`, http.StatusNotFound)
		return
	}

	ctx := r.Context()
	threads, err := h.threadSvc.ListAllThreads(ctx)
	if err != nil {
		logger.Error("failed to list all threads", "error", err)
		http.Error(w, `{"error": "Failed to list threads"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(threads); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}