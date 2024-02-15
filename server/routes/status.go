package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStatus reports if the server is operational.
//
// GET /status
func GetStatus(router *gin.RouterGroup) {
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Name":    "Bungo",
			"Version": "v0.0.1",
		})
	})
}
