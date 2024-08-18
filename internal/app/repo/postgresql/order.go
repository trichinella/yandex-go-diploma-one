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

	row := r.DB.QueryRow(context.Background(), "addOrder", order.ID, order.UserID, order.Number, order.CreatedDate, order.StatusID, order.Accrual, order.Paid)
	createdOrder := entity.Order{}
	if err := row.Scan(&createdOrder.ID, &createdOrder.UserID, &createdOrder.Number, &createdOrder.CreatedDate, &createdOrder.StatusID, &createdOrder.Accrual, &createdOrder.Paid); err != nil {
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
	if err := row.Scan(&foundOrder.ID, &foundOrder.UserID, &foundOrder.Number, &foundOrder.CreatedDate, &foundOrder.StatusID, &foundOrder.Accrual, &foundOrder.Paid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		logging.Sugar.Fatalf("Error search order by number: %v", err)
	}

	return &foundOrder, nil
}

func (r PostgresRepository) OrdersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
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

	rows, err := r.DB.Query(childCtx, "ordersByUser", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		order := entity.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.CreatedDate, &order.StatusID, &order.Accrual, &order.Paid); err != nil {
			logging.Sugar.Fatalf("Error search order by number: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r PostgresRepository) SaveOrder(ctx context.Context, order *entity.Order) error {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(childCtx, "saveOrder", `UPDATE public.orders 
SET user_id=$1, number=$2, created_date=$3, status_id=$4, accrual=$5, paid=$6
WHERE id=$7`)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(context.Background(), "saveOrder", order.UserID, order.Number, order.CreatedDate, order.StatusID, order.Accrual, order.Paid, order.ID)

	return err
}

func (r PostgresRepository) WithdrawalOrdersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(childCtx, "paidOrdersByUser", `SELECT id, user_id, number, created_date, status_id, accrual, paid 
FROM public.orders 
WHERE user_id = $1 AND paid > 0
ORDER BY created_date ASC
`)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(childCtx, "paidOrdersByUser", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		order := entity.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.CreatedDate, &order.StatusID, &order.Accrual, &order.Paid); err != nil {
			logging.Sugar.Fatalf("Error search order by number: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r PostgresRepository) NewOrders(ctx context.Context) ([]entity.Order, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(childCtx, "newOrders", `SELECT id, user_id, number, created_date, status_id, accrual, paid 
FROM public.orders 
WHERE status_id = $1
ORDER BY created_date ASC
`)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(childCtx, "newOrders", r.OrderStatusByCode(entity.NEW).ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		order := entity.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.CreatedDate, &order.StatusID, &order.Accrual, &order.Paid); err != nil {
			logging.Sugar.Fatalf("Error search order by number: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
