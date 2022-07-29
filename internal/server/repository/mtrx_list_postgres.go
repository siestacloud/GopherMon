package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/sirupsen/logrus"
)

type MtrxListPostgres struct {
	db *sqlx.DB
}

func NewMtrxListPostgres(db *sqlx.DB) *MtrxListPostgres {
	return &MtrxListPostgres{
		db: db,
	}
}

func (m *MtrxListPostgres) TestDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if m.db == nil {
		return errors.New("database are not connected")
	}
	return m.db.PingContext(ctx)
}

func (m *MtrxListPostgres) Create(mtrx *core.Metric) (int, error) {
	if m.db == nil {
		return 0, errors.New("database are not connected")
	}

	var mtrxID int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (name, type, value, delta) values ($1, $2,$3,$4) RETURNING id", mtrxTable)
	row := m.db.QueryRow(createItemQuery, mtrx.GetID(), mtrx.GetType(), mtrx.Value, mtrx.Delta)
	err := row.Scan(&mtrxID)
	if err != nil {
		return 0, err
	}
	return mtrxID, nil
}

// GetByNameType Получаю метрику из базы по имени и типу
func (m *MtrxListPostgres) Get(name string) (*core.Metric, error) {
	if m.db == nil {
		return nil, errors.New("database are not connected")
	}
	var dbMtrx core.Metric
	query := fmt.Sprintf(`SELECT name, type, value, delta FROM %s  WHERE name = $1`,
		mtrxTable)
	if err := m.db.Get(&dbMtrx, query, name); err != nil {
		return nil, err
	}
	return &dbMtrx, nil
}

func (m *MtrxListPostgres) Update(mtrx *core.Metric) (int, error) {
	if m.db == nil {
		return 0, errors.New("database are not connected")
	}
	var createItemQuery string

	if mtrx.GetType() == "counter" {
		mtrxVal, err := mtrx.GetDelta()
		if err != nil {
			return 0, err
		}
		createItemQuery = fmt.Sprintf("UPDATE %s SET delta = %v WHERE name = '%s' AND type = '%s' ", mtrxTable, mtrxVal, mtrx.GetID(), mtrx.GetType())

	} else {
		mtrxVal, err := mtrx.GetValue()
		if err != nil {
			return 0, err
		}
		createItemQuery = fmt.Sprintf("UPDATE %s SET value = %v WHERE name = '%s' AND type = '%s' ", mtrxTable, mtrxVal, mtrx.GetID(), mtrx.GetType())

	}

	_, err := m.db.Exec(createItemQuery)
	return 1, err
}

func (r *MtrxListPostgres) Flush(mtrxCase []core.Metric) (int, error) {

	// на всякий
	if r.db == nil {
		return 0, errors.New("You haven`t opened the database connection")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err

	}

	defer tx.Rollback()

	// stmt, err := tx.Prepare("INSERT INTO videos(title, description, views, likes) VALUES(?,?,?,?)")
	// stmt, err := tx.Prepare(`INSERT INTO mtrx (name, type)
	//                  VALUES($1,$2);`)

	// готовим инструкцию
	stmt, err := tx.Prepare(`INSERT INTO mtrx (name, type, value, delta) VALUES ($1, $2, $3, $4) ON CONFLICT (name) DO UPDATE SET type=$2, value = $3, delta = $4;`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, mtrx := range mtrxCase {
		if mtrx.MType == "gauge" {
			// указываю, что каждая метрика будет добавлена в транзакцию
			if _, err = stmt.Exec(
				// переменная в запросе
				mtrx.ID,
				mtrx.MType,
				*mtrx.Value,
				nil,
			); err != nil {
				if err = tx.Rollback(); err != nil {
					log.Fatalf("update drivers: unable to rollback: %v", err)
				}
				return 0, err
			}

		} else {
			// указываю, что каждая метрика будет добавлена в транзакцию
			if _, err = stmt.Exec(
				// переменная в запросе
				mtrx.ID,
				mtrx.MType,
				nil,
				*mtrx.Delta,
			); err != nil {
				if err = tx.Rollback(); err != nil {
					log.Fatalf("update drivers: unable to rollback: %v", err)
				}
				return 0, err
			}
		}

	}

	logrus.Warn("OK")
	//сохраняем изменения
	if err := tx.Commit(); err != nil {
		log.Fatalf("update drivers: unable to commit: %v", err)
	}

	return 1, nil
}
