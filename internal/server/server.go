package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/argo/signer"
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
	signer   *signer.Signer
	store    *store.Store
}

func New(database *database.Context, signer *signer.Signer, store *store.Store) *Context {
	engine := gin.Default()
	s := &Context{database: database, engine: engine, signer: signer, store: store}

	s.engine.Use(ErrorHandler)
	s.engine.GET("/status", s.Status)
	s.engine.POST("/data", s.UploadData)
	s.engine.POST("/data-item", s.UploadDataItem)

	return s
}

func (s *Context) Run(port string) {
	err := s.engine.Run(port)
	if err != nil {
		log.Fatalln("failed to start server", err)
	}
}
