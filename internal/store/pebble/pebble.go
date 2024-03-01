package pebble

import (
	"log"

	"github.com/cockroachdb/pebble"
)

type Database struct {
	db     *pebble.DB // Underlying pebble storage engine
	closed bool
}

func New(path string) *Database {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	return &Database{db: db, closed: false}
}

func (d *Database) Close() error {
	if d.closed {
		return nil
	}
	d.closed = true
	return d.db.Close()
}

func (d *Database) Delete(key []byte) error {
	return d.db.Delete(key, nil)
}

// Get retrieves the given key if it's present in the key-value store.
func (d *Database) Get(key []byte) ([]byte, error) {
	dat, closer, err := d.db.Get(key)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, len(dat))
	copy(ret, dat)
	err = closer.Close()
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// Has checks the given key if it's present in the key-value store.
func (d *Database) Has(key []byte) (bool, error) {
	_, closer, err := d.db.Get(key)
	if err != nil {
		return false, err
	}
	err = closer.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Put inserts the given value into the key-value store.
func (d *Database) Put(key []byte, value []byte) error {
	return d.db.Set(key, value, pebble.Sync)
}

