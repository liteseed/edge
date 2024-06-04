package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /status

func (s *Server) StatusGet(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := s.wallet.Client.GetInfo()
		if err != nil {
			c.JSON(
				http.StatusFailedDependency,
				gin.H{
					"Address": s.wallet.Signer.Address,
					"Name":    "Edge",
					"Version": version,
					"Gateway": gin.H{
						"Block-Height": "",
						"URL":          s.gateway,
						"Status":       "failed",
					},
				},
			)
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"Address": s.wallet.Signer.Address,
				"Name":    "Edge",
				"Version": version,
				"Gateway": gin.H{
					"Block-Height": info.Height,
					"URL":          s.gateway,
					"Status":       "ok",
				},
			},
		)
	}
}
