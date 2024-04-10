package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET Signer reports if the server is operational.
//
// GET /status

func getStatus(version string) gin.HandlerFunc   {
	return func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
			"Name": "Edge",
			"Version": version,
		})
	}
}
