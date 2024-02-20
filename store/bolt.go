package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	boltAllocSize = 8 * 1024 * 1024
	boltName      = "seed.db"
	BoltType      = "boltdb"
)

type BoltDB struct {
	Db *bolt.DB
}

func NewBoltDB(directory string) (*BoltDB, error) {
	if len(directory) == 0 {
		return nil, errors.New("bolt: directory can not null")
	}
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}

	Db, err := bolt.Open(path.Join(directory, boltName), 0660, &bolt.Options{Timeout: 2 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		if err == bolt.ErrTimeout {
			return nil, errors.New("cannot obtain database lock, database may be in use by another process")
		}
		return nil, err
	}
	Db.AllocSize = boltAllocSize
	boltDB := &BoltDB{
		Db: Db,
	}
	if err := boltDB.Db.Update(func(tx *bolt.Tx) error {
		buckets := []string{
			DataStore,
		}
		return createBuckets(tx, buckets)
	}); err != nil {
		return nil, err
	}
	return boltDB, nil
}

func (s *BoltDB) Type() string {
	return BoltType
}
func (s *BoltDB) Put(bucket, key string, value interface{}) (err error) {
	if _, ok := value.([]byte); !ok {
		return fmt.Errorf("unknown data type: %s, db: bolt db", reflect.TypeOf(value))
	}
	err = s.Db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		return bkt.Put([]byte(key), value.([]byte))
	})
	return
}

func (s *BoltDB) Get(bucket, key string) (data []byte, err error) {
	err = s.Db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket([]byte(bucket)).Get([]byte(key))
		if data == nil {
			return ErrDoesNotExist
		}
		return nil
	})
	return
}

func (s *BoltDB) GetAllKey(bucket string) (keys []string, err error) {
	keys = make([]string, 0)
	err = s.Db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})
	return
}

func (s *BoltDB) GetStream() error {
	return ErrUnimplemented
}

func (s *BoltDB) Delete(bucket, key string) (err error) {
	err = s.Db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete([]byte(key))
	})
	return
}

func (s *BoltDB) Exist(bucket, key string) bool {
	_, err := s.Get(bucket, key)
	return err == nil
}

func (s *BoltDB) Close() (err error) {
	return s.Db.Close()
}

func createBuckets(tx *bolt.Tx, buckets []string) error {
	for _, bucket := range buckets {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return err
		}
	}
	return nil
}

func NewBoltStore(directory string) (*Store, error) {
	Db, err := NewBoltDB(directory)
	if err != nil {
		return nil, err
	}
	log.Println("bolt connected - directory: " + directory)
	return &Store{KVDb: Db}, nil
}
