package store

import "github.com/liteseed/bungo/store/bolt"

func NewBoltStore(directory string) (*Store, error) {
	Db, err := bolt.NewBoltDB(directory)
	if err != nil {
		return nil, err
	}
	return &Store{KVDb: Db}, nil
}
