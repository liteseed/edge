package store

import (
	"os"

	"github.com/google/uuid"
)

type IStore interface {
	Put(bucket string, key string, value interface{}) (err error)

	Get(bucket string, key string) (data []byte, err error)

	GetStream(bucket string, key string) (data *os.File, err error)

	GetAllKey(bucket string) (keys []string, err error)

	Delete(bucket string, key string) (err error)

	Close() (err error)

	Type() string

	Exist(bucket, key string) bool
}

type Store struct {
	KVDb IStore
}

func (s *Store) Save(data []byte) (string, error) {
	id := uuid.New().String()
	s.KVDb.Put(data_store, id, data)
	return id, nil
}
