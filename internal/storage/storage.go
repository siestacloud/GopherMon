package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/siestacloud/service-monitoring/internal/mtrx"
)

type (
	//Хранит пул метрик
	Storage struct {
		filename string
		Mp       *mtrx.MetricsPool
	}
)

func NewStorage(filename string) (*Storage, error) {

	return &Storage{
		filename: filename,
		Mp:       mtrx.NewMetricsPool(),
	}, nil
}

func (s *Storage) TakeAll() ([]byte, error) {

	js, err := json.MarshalIndent(s.Mp, "", "	")
	if err != nil {
		return nil, err
	}

	return js, nil
}

func (s *Storage) WriteStorage() error {
	file, err := json.MarshalIndent(s.Mp, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.filename, file, 0644)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) ReadStorage() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0644) // открыть файл в режиме чтения если файла не существует, создать новый — флаг O_CREATE;
	if err != nil {
		return err
	}
	file.Close()
	mp, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(mp), s.Mp)
	if err != nil {
		s.Mp = mtrx.NewMetricsPool()
	}

	return nil
}
