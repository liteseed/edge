package api

import (
	"github.com/liteseed/bungo/cache"
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/store"
)

const MAX_DATA_ITEM_SIZE = 1_073_824

type API struct {
	cache *cache.Cache
	db    *database.Database
	store *store.Store
}

func New(
	c *cache.Cache,
	db *database.Database,
	s *store.Store,
) *API {
	return &API{cache: c, db: db, store: s}
}
