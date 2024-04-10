package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s := &Storage{db: db}
	return s, s.inits()
}

func MustInit(storagePath string) *Storage {
	storage, err := New(storagePath)
	if err != nil {
		panic("failed to init storage: " + err.Error())
	}
	return storage
}

func (s *Storage) inits() error {

	// create users table
	if err := s.userInit(); err != nil {
		return err
	}

	// create messages table
	if err := s.messageInit(); err != nil {
		return err
	}

	return nil
}
