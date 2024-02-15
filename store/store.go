package store

import "os"

type Store interface {
	Put(bucket, key string, value interface{}) (err error)

	Get(bucket, key string) (data []byte, err error)

	GetStream(bucket, key string) (data *os.File, err error)

	GetAllKey(bucket string) (keys []string, err error)

	Delete(bucket, key string) (err error)

	Close() (err error)

	Type() string

	Exist(bucket, key string) bool
}
