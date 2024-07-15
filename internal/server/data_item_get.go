package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DataItemGet Get /tx
func (srv *Server) DataItemGet(ctx *gin.Context) {
	id := ctx.Param("id")

	raw, err := srv.store.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "data item id does not exist"})
		return
	}

	ctx.Data(http.StatusOK, "application/octet-stream", raw)
}
