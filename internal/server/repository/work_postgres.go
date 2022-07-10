package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBPostgres struct {
	db *sqlx.DB
}

func newDBPostgres(db *sqlx.DB) *DBPostgres {
	return &DBPostgres{
		db: db,
	}
}

func (r *DBPostgres) TestDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if r.db == nil {
		return errors.New("database are not connected")
	}
	return r.db.PingContext(ctx)
}
