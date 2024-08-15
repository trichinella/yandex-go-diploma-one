package handler

import (
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/handler/body"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/algo/luhn"
	"diploma1/internal/app/service/balance"
	"diploma1/internal/app/service/logging"
	"errors"
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

func BalanceWithdrawHandle(userRepo repo.UserRepository, orderRepo repo.OrderRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := body.Content(r)
		if err != nil {
			logging.Sugar.Infof("error while getting content from request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = balance.Withdraw(r.Context(), userRepo, orderRepo, content)
		if err != nil {
			logging.Sugar.Error(err)
			var lackMoneyError *erroring.LackMoneyError
			if errors.As(err, &lackMoneyError) {
				w.WriteHeader(http.StatusPaymentRequired)
				return
			}

			var luhnNumberError *luhn.LuhnNumberError
			if errors.As(err, &luhnNumberError) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UserWithdrawalsHandle(orderRepo repo.OrderRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withdrawalOrders, err := balance.WithdrawalOrders(r.Context(), orderRepo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(withdrawalOrders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		rawBytes, err := easyjson.Marshal(withdrawalOrders)
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
