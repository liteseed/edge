package routes

import (
	"github.com/liteseed/argo/signer"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
)

const (
	CONTENT_TYPE_OCTET_STREAM = "application/octet-stream"
	MAX_DATA_SIZE             = 1_073_824
	MAX_DATA_ITEM_SIZE        = 1_073_824
)

type Routes struct {
	database *database.Database
	store    *store.Store
	signer   *signer.Signer
}

func New(
	database *database.Database,
	store *store.Store,
	signer *signer.Signer,
) *Routes {
	return &Routes{database: database, store: store, signer: signer}
}
