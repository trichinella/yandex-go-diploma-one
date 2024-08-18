package accrual

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/repo"
	"diploma1/internal/app/service/logging"
	"time"
)

func NewOrderQueue(orderRepo repo.OrderRepository, userRepo repo.UserRepository, orderChannel chan entity.Order) {
	for {
		select {
		case order := <-orderChannel:
			newOrderHandle(order, orderChannel, orderRepo, userRepo)
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func ReadOldOrders(orderRepo repo.OrderRepository, orderChannel chan<- entity.Order) {
	ctx := context.Background()
	orders, err := orderRepo.NewOrders(ctx)
	if err != nil {
		logging.Sugar.Fatal(err)
	}

	for _, order := range orders {
		go func() {
			orderChannel <- order
		}()
	}
}
