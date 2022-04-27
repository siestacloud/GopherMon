package mtrx

import (
	"encoding/json"
	"fmt"
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
func (m *MetricsPool) Add(key string, mtrx Metric) bool {
	if key == "" {
		return false
	}
	if m.LookUP(key) == nil {
		m.M[key] = mtrx
		return true //доб новую метрику в мапку
	}
	return false
}

//Обновить метрику в пуле
func (m *MetricsPool) Update(key string, mtrx Metric) bool {
	if key == "" {
		return false
	}
	switch mtrx.GetType() {
	case "counter":
		dmtrx := m.LookUP(key)
		if dmtrx == nil {
			return false
		}

		d, err := dmtrx.GetDelta()
		if err != nil {
			fmt.Println("1", err)
			return false
		}

		dm, err := mtrx.GetDelta()
		if err != nil {
			fmt.Println(err)
			return false
		}
		val := d + dm
		fmt.Println("HERE", val)

		err = mtrx.SetValue(val)
		if err != nil {
			return false
		}
	}
	m.M[key] = mtrx
	return true //доб новую метрику в мапку
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
		fmt.Printf("\n\nkey: %s value %v  \n", k, d) // вывести  все
	}
}

func (m *MetricsPool) GetAll() ([]byte, error) {

	js, err := json.MarshalIndent(m.M, "", "	")
	if err != nil {
		return nil, err
	}

	return js, nil
}

// func Save() error {
// 	fmt.Println("Saving", mtrxStorage)
// 	err := os.Remove(mtrxStorage)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	saveTo, err := os.Create(mtrxStorage)
// 	if err != nil {
// 		fmt.Println("Unable create file ", err)
// 		return err
// 	}
// 	defer saveTo.Close()

// 	encoder := gob.NewEncoder(saveTo)
// 	err = encoder.Encode(mtrxStorage)
// 	if err != nil {
// 		fmt.Println("Cannot save to ", mtrxStorage)
// 		return err
// 	}
// 	return nil
// }

// func Load() error {
// 	fmt.Println("Loading", mtrxStorage)
// 	loadFrom, err := os.Open(mtrxStorage)
// 	if err != nil {
// 		fmt.Println("Empty key/value store")
// 		return err
// 	}
// 	defer loadFrom.Close()
// 	decoder := gob.NewDecoder(loadFrom)
// 	err = decoder.Decode(&mtrxStorage)
// 	if err != nil {
// 		fmt.Println("Unable load from file ")
// 		return err
// 	}
// 	return nil
// }
