package postgresql

import (
	"context"
	"diploma1/internal/app/config"
	"diploma1/internal/app/service/logging"
	"fmt"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	DB *pgx.Conn
}

func GetPostgresRepository() *PostgresRepository {
	conn, err := pgx.Connect(context.Background(), config.State().DatabaseDSN)
	if err != nil {
		logging.Sugar.Fatal(fmt.Errorf("Unable to connect to database: %v\n", err))
	}
	pgxdecimal.Register(conn.TypeMap())

	return &PostgresRepository{
		DB: conn,
	}
}
