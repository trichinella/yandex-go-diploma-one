package postgresql

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/service/logging"
	"sync"
)

var once map[entity.StatusCode]entity.OrderStatus
var onceMutex sync.Mutex

func (r PostgresRepository) OrderStatus(statusCode entity.StatusCode) entity.OrderStatus {
	onceMutex.Lock()
	if once == nil {
		once = make(map[entity.StatusCode]entity.OrderStatus)
	}
	if _, ok := once[statusCode]; !ok {
		once[statusCode] = r.directOrderStatus(statusCode)
	}
	onceMutex.Unlock()

	return once[statusCode]
}

func (r PostgresRepository) directOrderStatus(statusCode entity.StatusCode) entity.OrderStatus {
	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(context.Background(), "status", `SELECT id, title, code FROM public.order_statuses WHERE code = $1`)
	if err != nil {
		logging.Sugar.Fatalf("Error prepare order status by code: %v", err)
	}

	row := r.DB.QueryRow(context.Background(), "status", statusCode)
	foundOrderStatus := entity.OrderStatus{}
	if err := row.Scan(&foundOrderStatus.ID, &foundOrderStatus.Title, &foundOrderStatus.Code); err != nil {
		logging.Sugar.Fatalf("Error search order status by code: %v", err)
	}

	return foundOrderStatus
}
