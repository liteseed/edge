package api

import (
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/store"
)

const MAX_DATA_ITEM_SIZE = 1_073_824

type API struct {
	database *database.Database
	store    *store.Store
}

func New(
	database *database.Database,
	store *store.Store,
) *API {
	return &API{database: database, store: store}
}
