package store

import (
	"github.com/google/uuid"
)

type IStore interface {
	Put(bucket string, key string, value interface{}) (err error)

	Get(bucket string, key string) (data []byte, err error)

	Delete(bucket string, key string) (err error)

	Close() (err error)

	Type() string

	Exist(bucket, key string) bool
}

type Store struct {
	KVDb IStore
}

func (s *Store) Save(data []byte) (uuid.UUID, error) {
	id := uuid.New()
	err := s.KVDb.Put(DataStore, id.String(), data)
	return id, err
}

func (s *Store) Get(id string) ([]byte, error) {
	return s.KVDb.Get(DataStore, id)
}
