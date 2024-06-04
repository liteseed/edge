package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get /tx
func (s *Server) DataItemGet(context *gin.Context) {
	id := context.Param("id")

	raw, err := s.store.Get(id)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusNotFound, gin.H{"error": "transaction id does not exist"})
		return
	}

	context.JSON(
		http.StatusOK,
		raw,
	)
}
