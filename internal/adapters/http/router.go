package http

import (
	"net/http"

	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/app/services"
)

func NewRouter(
	sessionSvc *services.SessionService,
	avatarSvc *services.AvatarService,
	threadSvc *services.ThreadService,
	commentSvc *services.CommentService,
	s3Client *s3.S3Client,
) http.Handler {
	mux := http.NewServeMux()
	sessionHandler := &SessionHandler{SessionService: sessionSvc}
	threadHandler := &ThreadHandler{
		threadSvc:  threadSvc,
		s3Client:   s3Client,
		sessionSvc: sessionSvc,
	}
	commentHandler := &CommentHandler{
		commentSvc: commentSvc,
		s3Client:   s3Client,
	}

	// API маршруты
	mux.HandleFunc("POST /api/session/name", sessionHandler.ChangeDisplayName)
	mux.HandleFunc("GET /api/session/me", sessionHandler.GetSessionInfo)
	mux.HandleFunc("GET /api/session/list", sessionHandler.ListSessions)

	mux.HandleFunc("POST /api/threads", threadHandler.CreateThread)
	mux.HandleFunc("GET /api/threads/", threadHandler.GetThread)
	mux.HandleFunc("GET /api/threads", threadHandler.ListActiveThreads)
	mux.HandleFunc("GET /api/threads/all", threadHandler.ListAllThreads)

	mux.HandleFunc("POST /api/threads/", commentHandler.CreateComment)
	mux.HandleFunc("GET /api/threads/", commentHandler.GetCommentsByThreadID)

	// HTML страницы
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		threads, err := threadSvc.ListActiveThreads(r.Context())
		if err != nil {
			RenderTemplate(w, "error.html", map[string]interface{}{
				"Code":    http.StatusInternalServerError,
				"Message": "Failed to load threads",
			})
			return
		}
		// Заполняем DisplayName для каждого треда
		for i, thread := range threads {
			session, err := sessionSvc.GetSessionByID(thread.SessionID.String())
			if err != nil {
				logger.Warn("failed to get session for thread", "thread_id", thread.ID, "error", err)
				threads[i].DisplayName = "Anonymous"
			} else {
				threads[i].DisplayName = session.DisplayName
			}
		}
		RenderTemplate(w, "catalog.html", threads)
	})

	mux.HandleFunc("GET /archive", func(w http.ResponseWriter, r *http.Request) {
		threads, err := threadSvc.ListAllThreads(r.Context())
		if err != nil {
			RenderTemplate(w, "error.html", map[string]interface{}{
				"Code":    http.StatusInternalServerError,
				"Message": "Failed to load archive",
			})
			return
		}
		// Заполняем DisplayName
		for i, thread := range threads {
			session, err := sessionSvc.GetSessionByID(thread.SessionID.String())
			if err != nil {
				logger.Warn("failed to get session for thread", "thread_id", thread.ID, "error", err)
				threads[i].DisplayName = "Anonymous"
			} else {
				threads[i].DisplayName = session.DisplayName
			}
		}
		RenderTemplate(w, "archive.html", threads)
	})

	mux.HandleFunc("GET /post/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/post/"):]
		id, err := utils.ParseUUID(idStr)
		if err != nil {
			RenderTemplate(w, "error.html", map[string]interface{}{
				"Code":    http.StatusBadRequest,
				"Message": "Invalid thread ID",
			})
			return
		}

		thread, err := threadSvc.GetThreadByID(r.Context(), id)
		if err != nil {
			RenderTemplate(w, "error.html", map[string]interface{}{
				"Code":    http.StatusNotFound,
				"Message": "Thread not found",
			})
			return
		}

		comments, err := commentSvc.GetCommentsByThreadID(r.Context(), id)
		if err != nil {
			RenderTemplate(w, "error.html", map[string]interface{}{
				"Code":    http.StatusInternalServerError,
				"Message": "Failed to load comments",
			})
			return
		}

		// Заполняем DisplayName и AvatarURL для треда
		session, err := sessionSvc.GetSessionByID(thread.SessionID.String())
		if err != nil {
			logger.Warn("failed to get session for thread", "thread_id", thread.ID, "error", err)
			thread.DisplayName = "Anonymous"
		} else {
			thread.DisplayName = session.DisplayName
		}

		// Заполняем DisplayName и AvatarURL для комментариев
		for i, comment := range comments {
			session, err := sessionSvc.GetSessionByID(comment.SessionID.String())
			if err != nil {
				logger.Warn("failed to get session for comment", "comment_id", comment.ID, "error", err)
				comments[i].DisplayName = "Anonymous"
				comments[i].AvatarURL = ""
			} else {
				comments[i].DisplayName = session.DisplayName
				comments[i].AvatarURL = session.AvatarURL
			}
		}

		isArchived := thread.IsDeleted
		data := map[string]interface{}{
			"Post":     thread,
			"Comments": comments,
		}

		if isArchived {
			RenderTemplate(w, "archive-post.html", data)
		} else {
			RenderTemplate(w, "post.html", data)
		}
	})

	mux.HandleFunc("GET /create-post", func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(w, "create-post.html", nil)
	})

	mux.HandleFunc("POST /submit-post", threadHandler.SubmitPost)
	mux.HandleFunc("POST /submit-comment", commentHandler.SubmitComment)

	// Статические файлы
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Middleware
	handler := SessionMiddleware(sessionSvc, "1337session")(mux)
	return handler
}