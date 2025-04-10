package http

import (
	"context"
	"net/http"

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
			var cookieVal string
			cookie, err := r.Cookie(cookieName)
			if err == nil {
				cookieVal = cookie.Value
			}

			sess, err := svc.GetOrCreate(cookieVal)
			if err != nil {
				http.Error(w, "failed to resolve session", http.StatusInternalServerError)
				return
			}

			// установить куку, если она ее не было или она была другой
			if cookie == nil || cookie.Value != sess.ID.String() {
				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    sess.ID.String(),
					Path:     "/",
					Expires:  sess.ExpiresAt,
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				})
			}

			ctx := context.WithValue(r.Context(), sessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
