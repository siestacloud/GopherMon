package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/siestacloud/service-monitoring/internal/mtrx"
)

type (
	//Хранит пул метрик
	Storage struct {
		Mp *mtrx.MetricsPool
		*W
		*R
	}
	//Открывает и записывает пул метрик в файл
	W struct {
		filename string
		file     *os.File // файл для записи
		// добавляем writer в Producer
		writer *bufio.Writer
	}
	//Читает из файла пул метрик
	R struct {
		file    *os.File // файл для чтения
		scanner *bufio.Scanner
	}
)

func NewStorage(filename string) (*Storage, error) {
	w, err := NewW(filename)
	if err != nil {
		return nil, err
	}
	r, err := NewR(filename)
	if err != nil {
		return nil, err
	}

	return &Storage{

		Mp: mtrx.NewMetricsPool(),
		W:  w,
		R:  r,
	}, nil
}

func NewW(filename string) (*W, error) {
	// открываем файл для записи в конец
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777) //открыть файл в режиме записи ли файла не существует, создать новый добавлять новые данные в файл
	if err != nil {
		return nil, err
	}

	return &W{file: file,
		filename: filename,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (w *W) WriteEvent(event *mtrx.MetricsPool) error {

	if err := os.Truncate(w.filename, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := w.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return w.writer.Flush()
}

func (w *W) Close() error {
	// закрываем файл
	return w.file.Close()
}

func NewR(filename string) (*R, error) {
	// открываем файл для чтения
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777) // открыть файл в режиме чтения если файла не существует, создать новый — флаг O_CREATE;

	if err != nil {
		return nil, err
	}

	return &R{file: file, // создаём новый Reader
		// reader: bufio.NewReader(file),
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil
}

func (r *R) Close() error {
	// закрываем файл
	return r.file.Close()
}
func (r *R) ReadEvent() (*mtrx.MetricsPool, error) {
	// // читаем данные до символа переноса строки
	// data, err := c.reader.ReadBytes('\n')
	// if err != nil {
	// 	return nil, err
	// }
	// одиночное сканирование до следующей строки
	if !r.scanner.Scan() {
		return nil, r.scanner.Err()
	}
	// читаем данные из scanner
	data := r.scanner.Bytes()

	// преобразуем данные из JSON-представления в структуру
	event := mtrx.MetricsPool{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (s *Storage) TakeAll() ([]byte, error) {

	js, err := json.MarshalIndent(s.Mp, "", "	")
	if err != nil {
		return nil, err
	}

	return js, nil
}
