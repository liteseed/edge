package server

import (
	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/internal/api"
)

type Server struct {
	engine *gin.Engine
}

func New() *Server {
	r := gin.Default()
	r.Use(ErrorHandler)

	return &Server{engine: r}
}

func (s *Server) Register(a *api.API) {
	s.engine.GET("/status", a.GetStatus)
	s.engine.GET("/:id", a.GetData)
	s.engine.POST("/", a.PostData)
}

func (s *Server) Run(port string) {
	s.engine.Run(port)
}
