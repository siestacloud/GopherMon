package service

import (
	"errors"

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

func (m *MtrxListService) Get(name string) (*core.Metric, error) {
	return m.repo.Get(name)
}

func (m *MtrxListService) Update(mtrx *core.Metric) (int, error) {
	return m.repo.Update(mtrx)
}

// Add обновление(создание) mtrx в базе
func (m *MtrxListService) Add(mtrx *core.Metric) (int, error) {
	dbMtrx, err := m.repo.Get(mtrx.ID)
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

	if mtrx.GetType() == dbMtrx.GetType() {
		if mtrx.GetType() == "counter" {
			sumDelta := *mtrx.Delta + *dbMtrx.Delta

			// сохраняю в базе
			err = mtrx.SetValue(sumDelta)
			if err != nil {
				return 0, err
			}
			logrus.Warn("mtrx delta = ", *mtrx.Delta)
		}

		id, err := m.repo.Update(mtrx)
		if err != nil {
			logrus.Warn("unable update  mtrx ", err)
			return 0, err
		}
		return id, nil

	}
	return 0, errors.New("mtrx in db have another type. drop this mtrx") //доб новую метрику в мапку
}

func (m *MtrxListService) Flush(mtrxCase []core.Metric) (int, error) {
	mtrxCaseOK := []core.Metric{}
	for _, mtrx := range mtrxCase {
		if mtrx.GetType() == "counter" {
			dbMtrx, err := m.repo.Get(mtrx.ID)
			if err != nil {
				logrus.Warn("mtrx not exist in postgres: ", err)
				mtrxCaseOK = append(mtrxCaseOK, mtrx)
				continue
			}
			if mtrx.GetType() == dbMtrx.GetType() {
				if mtrx.GetType() == "counter" {
					sumDelta := *mtrx.Delta + *dbMtrx.Delta
					// сохраняю в базе
					err = mtrx.SetValue(sumDelta)
					if err != nil {
						return 0, err
					}
				}
			}
		}
		mtrxCaseOK = append(mtrxCaseOK, mtrx)
	}
	return m.repo.Flush(mtrxCaseOK)
}
