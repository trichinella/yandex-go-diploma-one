package accrual

import (
	"context"
	"diploma1/internal/app/config"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/logging"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"strconv"
	"time"
)

type ReceivedAccrualOrder struct {
	Order   string            `json:"order"`
	Status  entity.StatusCode `json:"status"`
	Accrual float64           `json:"accrual,omitempty"`
}

type InitialOrder struct {
	Order string `json:"order"`
	Goods []struct {
		Description string `json:"description"`
		Price       int    `json:"price"`
	} `json:"goods"`
}

func newOrderHandle(order entity.Order, orderChannel chan entity.Order, orderRepo repo.OrderRepository, userRepo repo.UserRepository) {
	prepareAccrual(order)
	receivedAccrualOrder := ReceivedAccrualOrder{}
	path := fmt.Sprintf("http://%s/api/orders/%d", config.State().AccrualAddress, order.Number)
	logging.Sugar.Infow("Order has been prepared for checking",
		"Order number", order.Number)
	_, err := resty.New().R().SetResult(&receivedAccrualOrder).Get(path)
	if err != nil {
		logging.Sugar.Error(err)
		orderChannel <- order
		time.Sleep(3 * time.Second)
		return
	}

	logging.Sugar.Infow("Order checked",
		"Order number", receivedAccrualOrder.Order,
		"Accrual", receivedAccrualOrder.Accrual,
		"Status", receivedAccrualOrder.Status)

	if receivedAccrualOrder.Status == entity.NEW || receivedAccrualOrder.Status == entity.PROCESSING {
		orderChannel <- order
		time.Sleep(3 * time.Second)
		return
	}

	order.Accrual = receivedAccrualOrder.Accrual
	order.StatusId = orderRepo.OrderStatusByCode(receivedAccrualOrder.Status).ID
	//@todo mutex

	err = orderRepo.SaveOrder(context.Background(), &order)
	if err != nil {
		logging.Sugar.Error(err)
		orderChannel <- order
		time.Sleep(3 * time.Second)
		return
	}

	user, err := userRepo.UserById(context.Background(), order.UserId)
	if err != nil {
		logging.Sugar.Error(err)
		orderChannel <- order
		time.Sleep(3 * time.Second)
		return
	}

	user.Balance += receivedAccrualOrder.Accrual

	err = userRepo.SaveUser(context.Background(), user)
	if err != nil {
		logging.Sugar.Error(err)
		orderChannel <- order
		time.Sleep(3 * time.Second)
		return
	}
	//@todo mutex
}

// prepareAccrual этой функции быть не должно, но предоставленный бинарник работает не так, как заявлено
func prepareAccrual(order entity.Order) {
	orderNumber := strconv.Itoa(order.Number)
	initialOrder := InitialOrder{
		Order: orderNumber,
		Goods: []struct {
			Description string `json:"description"`
			Price       int    `json:"price"`
		}{
			{
				Description: "Чайник Bork",
				Price:       getPrice(),
			},
		},
	}

	path := fmt.Sprintf("http://%s/api/orders", config.State().AccrualAddress)
	_, err := resty.New().R().SetHeader("Content-Type", "application/json").SetBody(&initialOrder).Post(path)
	if err != nil {
		logging.Sugar.Error(err)
		return
	}
}

// getPrice с вероятностью 66% цена будет от 0, до 10000. С вероятностью 33% будет 0
func getPrice() int {
	if rand.Intn(3) == 1 {
		return 0
	}

	return rand.Intn(10000)
}
