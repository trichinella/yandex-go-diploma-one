package handler

import (
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/balance"
	"diploma1/internal/app/service/logging"
	"github.com/mailru/easyjson"
	"net/http"
)

func GetBalanceHandle(repository repo.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userFinanceState, err := balance.GetUserFinanceState(r.Context(), repository)

		if err != nil {
			logging.Sugar.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rawBytes, err := easyjson.Marshal(userFinanceState)
		if err != nil {
			logging.Sugar.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(rawBytes)

		if err != nil {
			logging.Sugar.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
