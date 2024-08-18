package middleware

import (
	"context"
	"diploma1/internal/app/service/auth"
	"diploma1/internal/app/service/ctxenv"
	"diploma1/internal/app/service/logging"
	"net/http"
)

func AuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerTokenString := r.Header.Get(`Authorization`)
			tokenString, err := auth.UnWrapTokenString(bearerTokenString)
			if err != nil {
				logging.Sugar.Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := auth.GetClaims(tokenString)
			if err != nil {
				logging.Sugar.Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxenv.ContextUserID, claims.UserID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
