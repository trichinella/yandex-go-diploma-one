package handler

import (
	"diploma1/internal/app/handler/body"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/auth"
	"diploma1/internal/app/service/logging"
	"errors"
	"net/http"
)

func RegisterHandle(repository repo.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := body.Content(r)
		if err != nil {
			logging.Sugar.Infof("error while getting content from request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tokenString, err := auth.Register(r.Context(), repository, content)
		if err != nil {
			logging.Sugar.Error(err)
			var loginExistsError *auth.LoginExistsError
			if errors.As(err, &loginExistsError) {
				w.WriteHeader(http.StatusConflict)
				return
			}

			if errors.Is(err, auth.ErrBadJson) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", auth.WrapTokenString(tokenString))
		w.WriteHeader(http.StatusOK)
	}
}
