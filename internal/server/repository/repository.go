package repository

import (
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

// Главный тип слоя repository, который встраивается как зависимость в слое SVC
type Repository struct {
	RAM
}

//Конструктор слоя repository
func NewRepository(mp *core.MetricsPool) *Repository {
	return &Repository{
		RAM: newRAMStorage(mp),
	}
}
