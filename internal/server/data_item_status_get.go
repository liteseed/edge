package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/edge/internal/database/schema"
)

// Get /tx/:id/status
func (s *Server) DataItemStatusGet(context *gin.Context) {
	id := context.Param("id")

	o, err := s.database.GetOrder(&schema.Order{ID: id})
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.String(http.StatusOK, string(o.Status))
}
