package postgresql

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/service/logging"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

func (r PostgresRepository) AddOrder(ctx context.Context, order entity.Order) (*entity.Order, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(childCtx, "addOrder", `INSERT INTO public.orders (id, user_id, number, created_date, status_id, accrual, paid) VALUES ($1, $2, $3, $4, $5, $6, $7)
	returning id,
	user_id,
	number,
	created_date,
	status_id,
	accrual,
	paid
`)
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(context.Background(), "addOrder", order.ID, order.UserId, order.Number, order.CreatedDate, order.StatusId, order.Accrual, order.Paid)
	createdOrder := entity.Order{}
	if err := row.Scan(&createdOrder.ID, &createdOrder.UserId, &createdOrder.Number, &createdOrder.CreatedDate, &createdOrder.StatusId, &createdOrder.Accrual, &createdOrder.Paid); err != nil {
		logging.Sugar.Fatalf("Error adding order: %v", err)
	}

	return &createdOrder, nil
}

func (r PostgresRepository) OrderByNumber(ctx context.Context, orderNumber int) (*entity.Order, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(childCtx, "orderByNumber", `SELECT id, user_id, number, created_date, status_id, accrual, paid FROM public.orders WHERE number = $1`)
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(childCtx, "orderByNumber", orderNumber)
	foundOrder := entity.Order{}
	if err := row.Scan(&foundOrder.ID, &foundOrder.UserId, &foundOrder.Number, &foundOrder.CreatedDate, &foundOrder.StatusId, &foundOrder.Accrual, &foundOrder.Paid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		logging.Sugar.Fatalf("Error search order by number: %v", err)
	}

	return &foundOrder, nil
}

func (r PostgresRepository) OrdersByUser(ctx context.Context, userId uuid.UUID) ([]entity.Order, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(childCtx, "ordersByUser", `SELECT id, user_id, number, created_date, status_id, accrual, paid 
FROM public.orders 
WHERE user_id = $1
ORDER BY created_date ASC
`)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(childCtx, "ordersByUser", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		order := entity.Order{}
		if err := rows.Scan(&order.ID, &order.UserId, &order.Number, &order.CreatedDate, &order.StatusId, &order.Accrual, &order.Paid); err != nil {
			logging.Sugar.Fatalf("Error search order by number: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
