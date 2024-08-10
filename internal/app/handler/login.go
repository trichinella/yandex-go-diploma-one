package handler

import (
	"diploma1/internal/app/handler/body"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/auth"
	"diploma1/internal/app/service/logging"
	"errors"
	"net/http"
)

func LoginHandle(repository repo.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := body.Content(r)
		if err != nil {
			logging.Sugar.Infof("error while getting content from request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tokenString, err := auth.Login(r.Context(), repository, content)
		if err != nil {
			logging.Sugar.Error(err)
			var loginAbsentError *auth.LoginAbsentError
			var incorrectPasswordError *auth.IncorrectPasswordError
			if errors.As(err, &loginAbsentError) || errors.As(err, &incorrectPasswordError) {
				w.WriteHeader(http.StatusUnauthorized)
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
