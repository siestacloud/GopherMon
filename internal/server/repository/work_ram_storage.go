package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/siestacloud/service-monitoring/internal/core"
)

type RAMStorage struct {
	db *core.MetricsPool
}

//NewMetricsPool Создает пул метрик
func newRAMStorage(mp *core.MetricsPool) *RAMStorage {
	return &RAMStorage{
		db: mp,
	}
}

func (r *RAMStorage) LookUP(key string) *core.Metric {
	return r.db.LookUP(key)
}

func (r *RAMStorage) Update(key string, mtrx *core.Metric) error {
	return r.db.Update(key, *mtrx)
}

func (r *RAMStorage) Create(key string, mtrx *core.Metric) error {
	return r.db.Create(key, *mtrx)
}

func (r *RAMStorage) PrintMtrxs() {
	r.db.PrintAll()
}

func (r *RAMStorage) GetAlljson() ([]byte, error) {
	js, err := json.MarshalIndent(r.db, "", "	")
	if err != nil {
		return nil, err
	}
	return js, nil
}

func (r *RAMStorage) WriteLocalStorage(fn string) error {
	file, err := json.MarshalIndent(r.db, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fn, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *RAMStorage) ReadLocalStorage(fn string) (*core.MetricsPool, error) {
	file, err := os.OpenFile(fn, os.O_RDONLY|os.O_CREATE, 0644) // открыть файл в режиме чтения если файла не существует, создать новый — флаг O_CREATE;
	if err != nil {
		return nil, err
	}
	file.Close()
	localMp, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	mp := core.NewMetricsPool()

	err = json.Unmarshal(localMp, mp)
	if err != nil {
		return mp, nil
	}
	return mp, nil
}
