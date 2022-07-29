package repository

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	mtrxTable = "mtrx"
)

func NewPostgresDB(urlDB string) (*sqlx.DB, error) {
	if urlDB == "" {
		return nil, errors.New("url not set")
	}
	db, err := sqlx.Open("postgres", urlDB)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Info("Success connect to postgres.")

	// делаем запрос
	var checkExist bool
	row := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  tablename  = 'mtrx');")
	err = row.Scan(&checkExist)
	if err != nil {
		log.Fatal(err)
	}
	if !checkExist {
		_, err = db.Exec("CREATE TABLE mtrx (id serial not null unique,name varchar(255) not null unique,type varchar(255) not null,value varchar(255), delta varchar(255) );") //QueryRowContext т.к. одна запись
		if err != nil {
			log.Fatal(err)
		}
		logrus.Info("Table mtrx successful create")

	} else {
		logrus.Info("Table mtrx already created")
	}

	return db, nil
}

// "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
