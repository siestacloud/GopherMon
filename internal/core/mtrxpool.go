package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

//CustomMetrics Пул метрик
type MetricsPool struct {
	M map[string]Metric
}

//NewMetricsPool Создает пул метрик
func NewMetricsPool() *MetricsPool {
	return &MetricsPool{
		M: map[string]Metric{},
	}
}

//Поиск метрики в общем пуле
func (m *MetricsPool) LookUP(key string) *Metric {
	_, ok := m.M[key]
	if ok {
		mtrx := m.M[key]
		return &mtrx
	} else {
		return nil
	}
}

//Добавление новой метрики
func (m *MetricsPool) Create(key string, mtrx Metric) error {
	if key == "" {
		return errors.New("unable create mtrx: empty key")
	}
	if m.LookUP(key) == nil {
		m.M[key] = mtrx
		return nil //доб новую метрику в мапку
	}
	return errors.New("unable create mtrx: mtrx already exist")
}

//Обновить метрику в пуле
func (m *MetricsPool) Update(key string, mtrx Metric) error {
	if key == "" {
		return errors.New("unable update mtrx: empty key")
	}
	switch mtrx.GetType() { //Определяю тип новой метрики
	// значение у нов метрики с типом счетчик не заменяет значение в базе а добавляется к нему.
	case "counter":
		dmtrx := m.LookUP(key) // ищу метрику по ключу
		if dmtrx == nil {
			break // метрики нету выход
		}

		d, err := dmtrx.GetDelta() // получаю значение метрики из базы
		if err != nil {
			return err
		}

		dm, err := mtrx.GetDelta() // получаю значение новой метрики
		if err != nil {
			return err
		}
		// суммирую значение
		val := d + dm
		// сохраняю в базе
		err = mtrx.SetValue(val)
		if err != nil {
			return err
		}
	}
	m.M[key] = mtrx // добавляю ее по ключу

	return nil //доб новую метрику в мапку
}

//Удалить метрику из пула
func (m *MetricsPool) Delete(key string) bool {
	if m.LookUP(key) != nil {
		delete(m.M, key)
		return true //удал метрику из мапки
	}
	return false
}

//Показать все метрики
func (m *MetricsPool) PrintAll() {
	for k, d := range m.M {
		infoPrint("mtrx from db", fmt.Sprintf("		MTRX: %s VAL: %v  \n", k, d))
	}
}

func (m *MetricsPool) GetAll() ([]byte, error) {

	js, err := json.MarshalIndent(m.M, "", "	")
	if err != nil {
		return nil, err
	}

	return js, nil
}

//WLS Write Local Storage
func (m *MetricsPool) WLS(fn string) error {
	file, err := json.MarshalIndent(m.M, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fn, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

// RLS Read Local Storage
func (m *MetricsPool) RLS(fn string) error {
	file, err := os.OpenFile(fn, os.O_RDONLY|os.O_CREATE, 0644) // открыть файл в режиме чтения если файла не существует, создать новый — флаг O_CREATE;
	if err != nil {
		return err
	}
	file.Close()
	localMp, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	// mp := NewMetricsPool()
	err = json.Unmarshal(localMp, m)
	if err != nil {
		return nil
	}
	return nil
}

func infoPrint(status, message string) {
	logrus.WithFields(logrus.Fields{"Layer": "repository", "Status": status}).Info(message)
}
