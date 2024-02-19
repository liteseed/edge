package cache

type Cache struct {
	Cache ICache
}

type ICache interface {
	Set(key string, entry []byte) error

	Get(key string) ([]byte, error)
}
