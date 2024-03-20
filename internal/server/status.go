package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET Status reports if the server is operational.
//
// GET /status
func (s *Context) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Name":    "Edge",
		"Version": "v0.0.1",
		"Signer":  s.signer.S.Address,
	})
}
