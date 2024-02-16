package server

import (
	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/cache"
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/store"
)

type Server struct {
	engine *gin.Engine
	cache  *cache.Cache
	db     *database.Database
	store  *store.Store
}

func New(c *cache.Cache, db *database.Database, s *store.Store) *Server {
	r := gin.Default()
	return &Server{engine: r, cache: c, db: db, store: s}
}

func (s *Server) Run(port string) {
	s.registerRoutes()
	s.engine.Run(port)
}

func (s *Server) registerRoutes() {
}
