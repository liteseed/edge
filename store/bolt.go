package store

import (
	"log"

	"github.com/liteseed/bungo/store/bolt"
)

func NewBoltStore(directory string) (*Store, error) {
	Db, err := bolt.NewBoltDB(directory)
	if err != nil {
		return nil, err
	}
	log.Println("bolt connected - directory: " + directory)
	return &Store{KVDb: Db}, nil
}
