package routes

import (
	"github.com/liteseed/bungo/internal/database"
	"github.com/liteseed/bungo/internal/store"
)

const MAX_DATA_ITEM_SIZE = 1_073_824

type Routes struct {
	database *database.Database
	store    *store.Store
}

func New(
	database *database.Database,
	store *store.Store,
) *Routes {
	return &Routes{database: database, store: store}
}
