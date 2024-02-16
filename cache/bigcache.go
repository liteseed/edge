package cache

import (
	"context"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

func NewBigCache(eviction time.Duration) (*Cache, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(eviction))
	if err != nil {
		return nil, err
	}

	log.Println("cache connected - eviction: " + eviction.String())
	return &Cache{Cache: cache}, nil
}

func (s *Cache) Set(key string, entry []byte) (err error) {
	return s.Cache.Set(key, entry)
}

func (s *Cache) Get(key string) ([]byte, error) {
	return s.Cache.Get(key)
}
