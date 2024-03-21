package server

import (
	"log"

	"github.com/everFinance/goar"
	"github.com/gin-gonic/gin"

	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
)

const (
	CONTENT_TYPE_OCTET_STREAM = "application/octet-stream"
	MAX_DATA_SIZE             = 1_073_824
	MAX_DATA_ITEM_SIZE        = 1_073_824
)

type Context struct {
	database *database.Context
	engine   *gin.Engine
	signer   *goar.Signer
	store    *store.Store
}

func New(database *database.Context, signer *goar.Signer, store *store.Store) *Context {
	engine := gin.New()

	engine.Use(gin.Recovery())
	s := &Context{database: database, engine: engine, signer: signer, store: store}

	s.engine.Use(ErrorHandler)
	s.engine.GET("/status", s.Status)
	s.engine.POST("/data/:id", s.uploadData)
	s.engine.POST("/data-item/:id", s.uploadDataItem)

	return s
}

func (s *Context) Run(port string) {
	err := s.engine.Run(port)
	if err != nil {
		log.Fatalln("failed to start server", err)
	}
}
