package store

import (
	"log"

	"github.com/liteseed/edge/internal/store/pebble"
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

func (s *Store) Delete(id string) error {
	return s.store.Delete([]byte(id))
}

func (s *Store) Get(id string) ([]byte, error) {
	return s.store.Get([]byte(id))
}

func (s *Store) Has(id string) (bool, error) {
	return s.store.Has([]byte(id))
}

func (s *Store) Put(id string, data []byte) error {
	err := s.store.Put([]byte(id), data)
	return err
}
