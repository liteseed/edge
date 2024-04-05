package server

import (
	"log/slog"
	"os"

	"github.com/everFinance/goar"
	"github.com/gin-gonic/gin"

	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
)

const (
	CONTENT_TYPE_OCTET_STREAM = "application/octet-stream"
	MAX_DATA_SIZE             = 1_073_824
	MAX_DATA_ITEM_SIZE        = 1_073_824
)

type Context struct {
	contract *contracts.Context
	database *database.Context
	engine   *gin.Engine
	signer   *goar.Signer
	store    *store.Store
	logger   *slog.Logger
}

func New(contract *contracts.Context, database *database.Context, logger *slog.Logger, store *store.Store) *Context {
	engine := gin.New()

	engine.Use(gin.Recovery())
	s := &Context{contract: contract, database: database, engine: engine, store: store}

	s.engine.Use(ErrorHandler)
	s.engine.GET("/status", s.Status)
	s.engine.POST("/data-item", s.uploadDataItem)

	return s
}

func (s *Context) Run(port string) {
	err := s.engine.Run(port)
	if err != nil {
		s.logger.Error("failed: server start", err)
		os.Exit(1)
	}
}
