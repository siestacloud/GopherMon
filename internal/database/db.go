package database

import (
	"os"
)

type Database struct {
	*os.File
}

func New(fn string) (*Database, error) {
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &Database{
		f,
	}, nil
}

func (s *Database) ReadMetrics() ([]byte, error) {
	buf, err := os.ReadFile(s.Name())
	if err != nil {
		return nil, err
	}
	return buf, nil
}
