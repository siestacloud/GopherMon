package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
// usersTable      = "users"
// todoListsTable  = "todo_lists"
// usersListsTable = "users_lists"
// todoItemsTable  = "todo_items"
// listsItemsTable = "lists_items"
)

func NewPostgresDB(urlDB string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", urlDB)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Info("Success connect to postgres.")
	return db, nil
}

// "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"


CREATE TABLE mtrx
(
    id      serial       not null unique,
    name    varchar(255) not null unique,
    type 	varchar(255) not null unique,
	value   varchar(255) not null unique,
    delta	varchar(255) not null
);