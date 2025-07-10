package database

import (
	"errors"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"mandacode.com/accounts/auth/ent"
)

func NewEntClient(dsn string) (*ent.Client, error) {
	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, errors.New("failed to connect to database: " + err.Error())
	}
	if drv == nil {
		return nil, errors.New("database driver is nil")
	}

	db := drv.DB()
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to ping database: " + err.Error())
	}

	client := ent.NewClient(ent.Driver(drv))
	if client == nil {
		return nil, errors.New("failed to create ent client")
	}

	return client, nil
}
