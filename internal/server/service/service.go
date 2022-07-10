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

type DB interface {
	TestDB() error
}

// Главный тип слоя SVC, который встраивается как зависимость в слое TRANSPORT
type Service struct {
	RAM
	DB
}

// Конструктор слоя SVC
func NewService(repos *repository.Repository) *Service {
	return &Service{
		RAM: newRAMService(repos.RAM),
		DB:  newDBService(repos.DB),
	}
}
