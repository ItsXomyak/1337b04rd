package http

import (
	"net/http"

	"1337b04rd/internal/app/services"
)

func NewRouter(
	sessionSvc *services.SessionService,
	avatarSvc *services.AvatarService,
) http.Handler {
	mux := http.NewServeMux()
	sessionHandler := &SessionHandler{SessionService: sessionSvc}

	mux.HandleFunc("POST /session/name", sessionHandler.ChangeDisplayName)
	mux.HandleFunc("GET /session/me", sessionHandler.GetSessionInfo)
	mux.HandleFunc("GET /session/list", sessionHandler.ListSessions)

	// middleware
	handler := SessionMiddleware(sessionSvc, "1337session")(mux)

	return handler
}
