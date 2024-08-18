package start

import (
	"diploma1/internal/app/config"
	"diploma1/internal/app/service/logging"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"strings"
)

func ExecMigrations() {
	m, err := migrate.New(
		"file://internal/migrations",
		strings.Replace(config.State().DatabaseDSN, "postgres", "pgx5", 1))
	if err != nil {
		logging.Sugar.Fatal(err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logging.Sugar.Fatal(err)
	}
}
