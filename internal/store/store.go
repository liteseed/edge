package store

import (
	"log"

	"github.com/google/uuid"
	"github.com/liteseed/bungo/internal/store/pebble"
)

type IStore interface {
	Close() (err error)

	Delete(key []byte) (err error)

	Get(key []byte) (data []byte, err error)

	Has(key []byte) (bool, error)

	Put(key []byte, value []byte) (err error)
}

type Store struct {
	store IStore
}

func New(storeOption string, path string) *Store {
	s := &Store{}
	switch storeOption {
	default:
		s.store = pebble.New(path)
	}
	return s
}

func (s *Store) Close() {
	err := s.store.Close()
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *Store) Delete(id uuid.UUID) error {
	return s.store.Delete([]byte(id.String()))
}

func (s *Store) Get(id uuid.UUID) ([]byte, error) {
	return s.store.Get([]byte(id.String()))
}

func (s *Store) Has(id uuid.UUID) (bool, error) {
	return s.store.Has([]byte(id.String()))
}

func (s *Store) Put(data []byte) (uuid.UUID, error) {
	id := uuid.New()
	err := s.store.Put([]byte(id.String()), data)
	return id, err
}
