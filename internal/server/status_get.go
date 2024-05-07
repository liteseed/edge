package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /status

func (s *Server) StatusGet(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"Name":    "Edge",
				"Version": version,
				"Address": s.wallet.Signer.Address,
			},
		)
	}
}
