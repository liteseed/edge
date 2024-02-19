package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStatus reports if the server is operational.
//
// GET /status
func (a *API) GetStatus(c *gin.Context) {
	log.Println("GET STATUS")
	c.JSON(http.StatusOK, gin.H{
		"Name":    "Bungo",
		"Version": "v0.0.1",
	})
}
