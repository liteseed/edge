package cache

import (
	"github.com/everFinance/goar/types"
	"github.com/liteseed/bungo/schema"
)

type Cache struct {
	Cache ICache
}

type ICache interface {
	Set(key string, entry []byte) error

	Get(key string) ([]byte, error)
}

func (c *Cache) GetFee() schema.ArFee {
	return schema.ArFee{Base: 1, PerChunk: 1}
}

func (c *Cache) GetInfo() types.NetworkInfo {
	return types.NetworkInfo{Height: 100}
}
