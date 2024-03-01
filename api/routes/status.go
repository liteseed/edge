package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET Status reports if the server is operational.
//
// GET /status
func (api *Routes) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Name":    "Bungo",
		"Version": "v0.0.1",
	})
}
