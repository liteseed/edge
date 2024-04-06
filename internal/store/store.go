package store

import (
	"log"

	"github.com/dgraph-io/badger"
)

type Store struct {
	store *badger.DB
}

func New(path string) *Store {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
	return &Store{store: db}
}

func (s *Store) Get(id string) ([]byte, error) {
	var ival []byte
	err := s.store.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		ival, err = item.ValueCopy(nil)
		return err
	})

	return ival, err
}

func (s *Store) Set(id string, data []byte) error {
	err := s.store.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(id), data)
	})
	return err
}

func (s *Store) Shutdown() error {
	return s.store.Close()
}
