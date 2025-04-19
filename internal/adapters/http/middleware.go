package http

import (
	"context"
	"net/http"

	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/app/services"
	"1337b04rd/internal/domain/session"
)

type contextKey string

const sessionKey contextKey = "session"

func GetSessionFromContext(ctx context.Context) (*session.Session, bool) {
	sess, ok := ctx.Value(sessionKey).(*session.Session)
	return sess, ok
}

func SessionMiddleware(svc *services.SessionService, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var sess *session.Session
			var sessionErr error

			cookie, err := r.Cookie(cookieName)
			if err == nil {
				sess, sessionErr = svc.GetOrCreate(ctx, cookie.Value)
			}

			if sessionErr != nil || sess == nil {
				logger.Warn("creating new session", "reason", sessionErr)
				sess, err = svc.CreateNew(ctx)
				if err != nil {
					logger.Error("failed to create new session", "error", err)
					Respond(w, http.StatusInternalServerError, map[string]string{"error": "failed to create session"})
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    sess.ID.String(),
					Path:     "/",
					Expires:  sess.ExpiresAt,
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					Secure:   false,
				})
				logger.Info("set new session cookie", "session_id", sess.ID)
			}

			ctx = context.WithValue(ctx, sessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
