package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/services"
	"1337b04rd/internal/domain/errors"
)

type CommentHandler struct {
	commentSvc *services.CommentService
	s3Client   *s3.S3Client
}

func NewCommentHandler(commentSvc *services.CommentService, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{
		commentSvc: commentSvc,
	}
}

func (h *CommentHandler) SubmitComment(w http.ResponseWriter, r *http.Request) {
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

    threadIDStr := r.FormValue("thread_id")
    content := r.FormValue("comment")
    parentIDStr := r.FormValue("parent_id")
    file, fileHeader, err := r.FormFile("file")
    var imageURL *string
    if err == nil && file != nil && fileHeader != nil {
        defer file.Close()
        contentType := fileHeader.Header.Get("Content-Type")
        url, err := h.s3Client.UploadCommentImage(file, fileHeader.Filename, contentType)
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

    threadID, err := utils.ParseUUID(threadIDStr)
    if err != nil {
        logger.Warn("invalid thread ID", "thread_id", threadIDStr)
        RenderTemplate(w, "error.html", map[string]interface{}{
            "Code":    http.StatusBadRequest,
            "Message": "Invalid thread ID",
        })
        return
    }

    var parentID *utils.UUID
    if parentIDStr != "" {
        pid, err := utils.ParseUUID(parentIDStr)
        if err != nil {
            logger.Warn("invalid parent ID", "parent_id", parentIDStr)
            RenderTemplate(w, "error.html", map[string]interface{}{
                "Code":    http.StatusBadRequest,
                "Message": "Invalid parent ID",
            })
            return
        }
        parentID = &pid
    }

    _, err = h.commentSvc.CreateComment(r.Context(), threadID, parentID, content, imageURL, sess.ID)
    if err != nil {
        logger.Error("failed to create comment", "error", err)
        RenderTemplate(w, "error.html", map[string]interface{}{
            "Code":    http.StatusInternalServerError,
            "Message": "Failed to create comment",
        })
        return
    }

    http.Redirect(w, r, "/post/"+threadID.String(), http.StatusSeeOther)
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	path := r.URL.Path
	if !strings.HasPrefix(path, "/threads/") || !strings.Contains(path, "/comments") {
		logger.Warn("invalid path", "path", path)
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}

	//	/threads/{thread_id}/comments
	parts := strings.Split(strings.TrimPrefix(path, "/threads/"), "/")
	if len(parts) < 2 || parts[1] != "comments" {
		logger.Warn("invalid path format", "path", path)
		http.Error(w, `{"error": "Invalid thread ID"}`, http.StatusBadRequest)
		return
	}
	threadIDStr := parts[0]
	threadID, err := utils.ParseUUID(threadIDStr)
	if err != nil {
		logger.Error("invalid thread_id", "error", err, "thread_id", threadIDStr)
		http.Error(w, `{"error": "Invalid thread ID"}`, http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		logger.Error("failed to parse form", "error", err)
		http.Error(w, `{"error": "Invalid form data"}`, http.StatusBadRequest)
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
			http.Error(w, `{"error": "Invalid parent_id"}`, http.StatusBadRequest)
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
		http.Error(w, `{"error": "Invalid session_id"}`, http.StatusBadRequest)
		return
	}

	if content == "" {
		logger.Warn("missing content", "thread_id", threadID)
		http.Error(w, `{"error": "Content is required"}`, http.StatusBadRequest)
		return
	}

	comment, err := h.commentSvc.CreateComment(ctx, threadID, parentID, content, imageURLPtr, sessionID)
	if err != nil {
		if err == errors.ErrThreadNotFound {
			logger.Warn("thread not found", "thread_id", threadID)
			http.Error(w, `{"error": "Thread not found"}`, http.StatusNotFound)
			return
		}
		logger.Error("failed to create comment", "error", err, "thread_id", threadID)
		http.Error(w, `{"error": "Failed to create comment"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

func (h *CommentHandler) GetCommentsByThreadID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warn("invalid method", "method", r.Method, "path", r.URL.Path)
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	path := r.URL.Path
	if !strings.HasPrefix(path, "/threads/") || !strings.Contains(path, "/comments") {
		logger.Warn("invalid path", "path", path)
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}

	// Извлекаем thread_id из пути: /threads/{thread_id}/comments
	parts := strings.Split(strings.TrimPrefix(path, "/threads/"), "/")
	if len(parts) < 2 || parts[1] != "comments" {
		logger.Warn("invalid path format", "path", path)
		http.Error(w, `{"error": "Invalid thread ID"}`, http.StatusBadRequest)
		return
	}
	threadIDStr := parts[0]
	threadID, err := utils.ParseUUID(threadIDStr)
	if err != nil {
		logger.Error("invalid thread_id", "error", err, "thread_id", threadIDStr)
		http.Error(w, `{"error": "Invalid thread ID"}`, http.StatusBadRequest)
		return
	}

	comments, err := h.commentSvc.GetCommentsByThreadID(ctx, threadID)
	if err != nil {
		if err == errors.ErrThreadNotFound {
			logger.Warn("thread not found", "thread_id", threadID)
			http.Error(w, `{"error": "Thread not found"}`, http.StatusNotFound)
			return
		}
		logger.Error("failed to get comments", "error", err, "thread_id", threadID)
		http.Error(w, `{"error": "Failed to get comments"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}