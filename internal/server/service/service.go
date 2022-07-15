package service

import (
	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/siestacloud/service-monitoring/internal/server/repository"
)

// интерфейс для работы со слоем Service
type RAM interface {
	GetAlljson() ([]byte, error)
	LookUP(key string) *core.Metric
	Add(key string, mtrx *core.Metric) error
	Create(key string, mtrx *core.Metric) error
	Update(key string, mtrx *core.Metric) error
	ReadLocalStorage(fn string) (*core.MetricsPool, error)
	WriteLocalStorage(fn string) error
	CheckHash(key string, mtrx *core.Metric) error
}

type MtrxList interface {
	TestDB() error
	Create(mtrx *core.Metric) (int, error)
	Get(name, mtype string) (*core.Metric, error)
	Update(mtrx *core.Metric) (int, error)
	Add(mtrx *core.Metric) (int, error)
}

// Главный тип слоя SVC, который встраивается как зависимость в слое TRANSPORT
type Service struct {
	RAM
	MtrxList
}

// Конструктор слоя SVC
func NewService(repos *repository.Repository) *Service {
	return &Service{
		RAM:      newRAMService(repos.RAM),
		MtrxList: NewMtrxListService(repos.MtrxList),
	}
}
