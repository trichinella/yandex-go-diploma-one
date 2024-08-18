package server

import (
	"context"
	"diploma1/internal/app/config"
	"diploma1/internal/app/entity"
	"diploma1/internal/app/handler"
	"diploma1/internal/app/middleware"
	"diploma1/internal/app/repo/postgresql"
	"diploma1/internal/app/service/accrual"
	"diploma1/internal/app/service/logging"
	"diploma1/internal/app/service/start"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	Router *chi.Mux
	Repo   *postgresql.PostgresRepository
}

func (s *Server) Run() {
	defer func() {
		_ = logging.Sugar.Sync()
	}()

	//Инициализация
	//Миграция
	start.ExecMigrations()
	//Accrual сервис проинициализировать стартовыми данными
	start.AccrualFilling()

	orderChannel := make(chan entity.Order, 100)

	s.Router.Group(func(r chi.Router) {
		r.Use(middleware.LogMiddleware())
		r.Use(middleware.Compress())

		r.Post(`/api/user/register`, handler.RegisterHandle(s.Repo))
		r.Post(`/api/user/login`, handler.LoginHandle(s.Repo))
	})

	s.Router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware())
		r.Use(middleware.LogMiddleware())
		r.Use(middleware.Compress())

		r.Post(`/api/user/orders`, handler.AddingOrderHandle(s.Repo, orderChannel))
		r.Get(`/api/user/orders`, handler.GetUserOrderListHandle(s.Repo))
		r.Get(`/api/user/balance`, handler.GetBalanceHandle(s.Repo))
		r.Post(`/api/user/balance/withdraw`, handler.BalanceWithdrawHandle(s.Repo, s.Repo))
		r.Get(`/api/user/withdrawals`, handler.UserWithdrawalsHandle(s.Repo))
	})

	go func() {
		accrual.ReadOldOrders(s.Repo, orderChannel)
	}()

	go func() {
		accrual.NewOrderQueue(s.Repo, s.Repo, orderChannel)
	}()

	logging.Sugar.Infow("Listen and serve", "Host", config.State().GopherMartAddress)
	err := http.ListenAndServe(config.State().GopherMartAddress, s.Router)
	if err != nil {
		logging.Sugar.Fatal(err)
	}
}

func Run() {
	server := newServer()
	server.Run()

	defer func() {
		err := server.Repo.DB.Close(context.Background())
		if err != nil {
			logging.Sugar.Fatal(fmt.Errorf("unable to close to database: %v", err))
		}
	}()
}

func newServer() *Server {
	return &Server{
		Repo:   postgresql.GetPostgresRepository(),
		Router: chi.NewRouter(),
	}
}
