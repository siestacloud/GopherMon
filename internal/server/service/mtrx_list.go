package service

import (
	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/siestacloud/service-monitoring/internal/server/repository"
	"github.com/sirupsen/logrus"
)

type MtrxListService struct {
	repo repository.MtrxList
}

func NewMtrxListService(repo repository.MtrxList) *MtrxListService {
	return &MtrxListService{
		repo: repo,
	}
}

func (m *MtrxListService) TestDB() error {
	return m.repo.TestDB()
}

func (m *MtrxListService) Create(mtrx *core.Metric) (int, error) {
	return m.repo.Create(mtrx)
}

func (m *MtrxListService) Get(name, mtype string) (*core.Metric, error) {
	return m.repo.Get(name, mtype)
}

func (m *MtrxListService) Update(mtrx *core.Metric) (int, error) {
	return m.repo.Update(mtrx)
}

func (m *MtrxListService) Add(mtrx *core.Metric) (int, error) {

	dbMtrx, err := m.repo.Get(mtrx.ID, mtrx.MType)
	if err != nil {
		logrus.Warn("mtrx not exist in postgres: ", err)
		logrus.Warn("create mtrx...")
		id, err := m.repo.Create(mtrx)
		if err != nil {
			logrus.Warn("unable create mtrx ", err)
			return 0, err
		}
		logrus.Warn("mtrx created")
		return id, nil
	}
	logrus.Warn("mtrx exist in db")
	logrus.Warn("update mtrx...")

	switch mtrx.MType {
	case "counter":
		logrus.Warn("mtrx type counter")

		sumDelta := *mtrx.Delta + *dbMtrx.Delta
		mtrx.Delta = &sumDelta
		logrus.Warn("mtrx delta = ", mtrx.Delta)

		id, err := m.repo.Update(mtrx)
		if err != nil {
			logrus.Warn("unable update mtrx ", err)
			return 0, err
		}
		logrus.Warn("mtrx updated")
		return id, err
	}

	return m.repo.Update(mtrx)
}
