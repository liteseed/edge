package server

import (
	"github.com/gin-gonic/gin"
	"github.com/liteseed/edge/api/routes"
)

type Server struct {
	engine *gin.Engine
}

func New() *Server {
	r := gin.Default()
	r.Use(ErrorHandler)

	return &Server{engine: r}
}

func (s *Server) Register(a *routes.Routes) {
	s.engine.GET("/status", a.Status)
	s.engine.POST("/data", a.UploadData)
	s.engine.POST("/dataitem", a.UploadDataItem)
}

func (s *Server) Run(port string) {
	s.engine.Run(port)
}
