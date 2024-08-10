package postgresql

import (
	"context"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/service/logging"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

func (r PostgresRepository) AddUser(ctx context.Context, user entity.User) (*entity.User, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(childCtx, "addUser", `INSERT INTO public.users (id, login, password, balance, created_date) VALUES ($1, $2, $3, $4, $5)
	returning id,
	login,
	password,
	balance,
	created_date
`)
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(context.Background(), "addUser", user.ID, user.Login, user.Password, user.Balance, user.CreatedDate)
	createdUser := entity.User{}
	if err := row.Scan(&createdUser.ID, &createdUser.Login, &createdUser.Password, &createdUser.Balance, &createdUser.CreatedDate); err != nil {
		logging.Sugar.Fatalf("Error adding user: %v", err)
	}

	return &createdUser, nil
}

func (r PostgresRepository) UserByLogin(ctx context.Context, login string) (*entity.User, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	//в pgx можно отдельно не готовить - внутри делает хэш
	_, err := r.DB.Prepare(context.Background(), "userByLogin", `SELECT id, login, password, balance, created_date FROM public.users WHERE login = $1`)
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(childCtx, "userByLogin", login)
	foundUser := entity.User{}
	if err := row.Scan(&foundUser.ID, &foundUser.Login, &foundUser.Password, &foundUser.Balance, &foundUser.CreatedDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		logging.Sugar.Fatalf("Error search user by login: %v", err)
	}

	return &foundUser, nil
}
