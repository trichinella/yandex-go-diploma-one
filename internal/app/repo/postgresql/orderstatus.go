package postgresql

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/service/logging"
	"github.com/google/uuid"
	"sync"
)

var once map[entity.StatusCode]entity.OrderStatus
var onceMutex sync.Mutex

var onceId map[uuid.UUID]entity.OrderStatus
var onceMutexId sync.Mutex

func (r PostgresRepository) OrderStatusByCode(statusCode entity.StatusCode) entity.OrderStatus {
	onceMutex.Lock()
	if once == nil {
		once = make(map[entity.StatusCode]entity.OrderStatus)
	}
	if _, ok := once[statusCode]; !ok {
		once[statusCode] = r.directOrderStatusByCode(statusCode)
	}
	onceMutex.Unlock()

	return once[statusCode]
}

func (r PostgresRepository) directOrderStatusByCode(statusCode entity.StatusCode) entity.OrderStatus {
	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(context.Background(), "statusCode", `SELECT id, title, code FROM public.order_statuses WHERE code = $1`)
	if err != nil {
		logging.Sugar.Fatalf("Error prepare order status by code: %v", err)
	}

	row := r.DB.QueryRow(context.Background(), "statusCode", statusCode)
	foundOrderStatus := entity.OrderStatus{}
	if err := row.Scan(&foundOrderStatus.ID, &foundOrderStatus.Title, &foundOrderStatus.Code); err != nil {
		logging.Sugar.Fatalf("Error search order status by code: %v", err)
	}

	return foundOrderStatus
}

func (r PostgresRepository) OrderStatusById(statusId uuid.UUID) entity.OrderStatus {
	onceMutexId.Lock()
	if onceId == nil {
		onceId = make(map[uuid.UUID]entity.OrderStatus)
	}
	if _, ok := onceId[statusId]; !ok {
		onceId[statusId] = r.directOrderStatusById(statusId)
	}
	onceMutexId.Unlock()

	return onceId[statusId]
}

func (r PostgresRepository) directOrderStatusById(statusId uuid.UUID) entity.OrderStatus {
	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(context.Background(), "statusId", `SELECT id, title, code FROM public.order_statuses WHERE id = $1`)
	if err != nil {
		logging.Sugar.Fatalf("Error prepare order status by code: %v", err)
	}

	row := r.DB.QueryRow(context.Background(), "statusId", statusId)
	foundOrderStatus := entity.OrderStatus{}
	if err := row.Scan(&foundOrderStatus.ID, &foundOrderStatus.Title, &foundOrderStatus.Code); err != nil {
		logging.Sugar.Fatalf("Error search order status by id: %v", err)
	}

	return foundOrderStatus
}
