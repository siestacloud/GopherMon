package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/service-monitoring/internal/core"
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

func (r *MtrxListPostgres) Create(mtrx *core.Metric) (int, error) {

	var mtrxId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (name, type, value, delta) values ($1, $2,$3,$4) RETURNING id", mtrxTable)
	row := r.db.QueryRow(createItemQuery, mtrx.GetID(), mtrx.GetType(), mtrx.Value, mtrx.Delta)
	err := row.Scan(&mtrxId)
	if err != nil {
		return 0, err
	}
	return mtrxId, nil
}

// GetByNameType Получаю метрику из базы по имени и типу
func (m *MtrxListPostgres) Get(name, mtype string) (*core.Metric, error) {
	var dbMtrx core.Metric
	query := fmt.Sprintf(`SELECT name, type, value, delta FROM %s  WHERE name = $1 AND type = $2`,
		mtrxTable)
	if err := m.db.Get(&dbMtrx, query, name, mtype); err != nil {
		return nil, err
	}
	return &dbMtrx, nil
}

func (r *MtrxListPostgres) Update(mtrx *core.Metric) (int, error) {
	var mtrxId int
	createItemQuery := fmt.Sprintf("UPDATE %s SET value=%d, delta=%d WHERE name = '%s', type = '%s' ", mtrxTable, mtrx.Value, mtrx.Delta, mtrx.GetID(), mtrx.GetType())
	row := r.db.QueryRow(createItemQuery)
	err := row.Scan(&mtrxId)
	if err != nil {
		return 0, err
	}
	return mtrxId, nil
}
