package http

import (
	"1337b04rd/internal/app/services"
	"net/http"
)

func NewRouter(
	sessionSvc *services.SessionService,
	avatarSvc *services.AvatarService,
	threadSvc *services.ThreadService,
	commentSvc *services.CommentService,
) http.Handler {
	mux := http.NewServeMux()
	sessionHandler := &SessionHandler{SessionService: sessionSvc}
	threadHandler := &ThreadHandler{threadSvc: threadSvc}
	commentHandler := &CommentHandler{commentSvc: commentSvc}

	// === Сессии ===
	mux.HandleFunc("POST /session/name", sessionHandler.ChangeDisplayName)
	mux.HandleFunc("GET /session/me", sessionHandler.GetSessionInfo)
	mux.HandleFunc("GET /session/list", sessionHandler.ListSessions)

	// === Треды ===
	mux.HandleFunc("POST /threads", threadHandler.CreateThread)
	mux.HandleFunc("GET /threads/view/", threadHandler.GetThread) // GET /threads/view/{id}
	mux.HandleFunc("GET /threads", threadHandler.ListActiveThreads)
	mux.HandleFunc("GET /threads/all", threadHandler.ListAllThreads)

	// === Комментарии ===
	mux.HandleFunc("POST /threads/comment", commentHandler.CreateComment)
	mux.HandleFunc("GET /threads/comment", commentHandler.GetCommentsByThreadID)

	// === Middleware ===
	handler := SessionMiddleware(sessionSvc, "1337session")(mux)

	return handler
}
