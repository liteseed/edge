package server

import (
	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/api"
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/store"
)

type Server struct {
	engine *gin.Engine
}

func New(db *database.Database, s *store.Store) *Server {
	a := api.New(db, s)

	r := gin.Default()

	r.GET("/status", a.GetStatus)

	r.GET("/:id", a.GetData)
	r.POST("/", a.PostData)

	return &Server{engine: r}
}

func (s *Server) Run(port string) {
	s.engine.Run(port)
}
