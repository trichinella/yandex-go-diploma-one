package handler

import (
	"diploma1/internal/app/entity"
	"diploma1/internal/app/erroring"
	"diploma1/internal/app/handler/body"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/algo/luhn"
	"diploma1/internal/app/service/logging"
	"diploma1/internal/app/service/order"
	"errors"
	"github.com/mailru/easyjson"
	"net/http"
)

func AddingOrderHandle(repository repo.OrderRepository, orderChannel chan<- entity.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := body.Content(r)
		if err != nil {
			logging.Sugar.Infof("error while getting content from request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = order.AddOrder(r.Context(), repository, content, orderChannel)
		if err != nil {
			var numberExistsError *erroring.NumberExistsError
			if errors.As(err, &numberExistsError) {
				w.WriteHeader(http.StatusOK)
				return
			}

			logging.Sugar.Error(err)
			if errors.Is(err, erroring.ErrIncorrectNumber) || errors.Is(err, erroring.ErrEmptyRequest) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var luhnNumberError *luhn.LuhnNumberError
			if errors.As(err, &luhnNumberError) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			var someoneElseOrderError *erroring.SomeoneElseOrderError
			if errors.As(err, &someoneElseOrderError) {
				w.WriteHeader(http.StatusConflict)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func GetUserOrderListHandle(repository repo.OrderRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userOrderList, err := order.GetUserOrderList(r.Context(), repository)
		if err != nil {
			logging.Sugar.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(userOrderList) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		rawBytes, err := easyjson.Marshal(userOrderList)
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
