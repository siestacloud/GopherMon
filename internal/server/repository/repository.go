package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/service-monitoring/internal/core"
)

// Интерфейс для взаимодействия с метриками, хранящимися в оперативе
type RAM interface {
	PrintMtrxs()
	GetAlljson() ([]byte, error)
	LookUP(key string) *core.Metric
	Create(key string, mtrx *core.Metric) error
	Update(key string, mtrx *core.Metric) error
	ReadLocalStorage(fn string) (*core.MetricsPool, error)
	WriteLocalStorage(fn string) error
}

// Интерфейс для взаимодействия с метриками, хранящимися в postgres
type MtrxList interface {
	TestDB() error
	Create(mtrx *core.Metric) (int, error)
	Get(name string) (*core.Metric, error)
	Update(mtrx *core.Metric) (int, error)
}

// Главный тип слоя repository, который встраивается как зависимость в слое SVC
type Repository struct {
	RAM
	MtrxList
}

//Конструктор слоя repository
func NewRepository(mp *core.MetricsPool, db *sqlx.DB) *Repository {
	return &Repository{
		RAM:      newRAMStorage(mp),
		MtrxList: NewMtrxListPostgres(db),
	}
}
