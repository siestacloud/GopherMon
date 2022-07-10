package service

import (
	"github.com/siestacloud/service-monitoring/internal/server/repository"
)

type DBService struct {
	repo repository.DB
}

func newDBService(repo repository.DB) *DBService {
	return &DBService{
		repo: repo,
	}
}

func (s *DBService) TestDB() error {
	return s.repo.TestDB()
}
