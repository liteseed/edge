package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get /tx
func (srv *Server) DataItemGet(ctx *gin.Context) {
	id := ctx.Param("id")

	raw, err := srv.store.Get(id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "data item id does not exist"})
		return
	}

	ctx.JSON(http.StatusOK, raw)
}
