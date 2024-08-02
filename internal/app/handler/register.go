package handler

import (
	"diploma1/internal/app/handler/body"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/auth"
	"errors"
	"net/http"
)

func RegisterHandle(repository repo.UserRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := body.Content(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = auth.Register(repository, content)
		if err != nil {
			if errors.Is(err, auth.ErrLoginExists) {
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

		w.WriteHeader(http.StatusOK)
	}
}
