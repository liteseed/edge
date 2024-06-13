package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /status

func (srv *Server) StatusGet(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := srv.client.GetNetworkInfo()
		if err != nil {
			c.JSON(
				http.StatusFailedDependency,
				gin.H{
					"Address": srv.signer.Address,
					"Name":    "Edge",
					"Version": version,
					"Gateway": gin.H{
						"Block-Height": "",
						"URL":          srv.client.Gateway,
						"Status":       "failed",
					},
				},
			)
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"Address": srv.signer.Address,
				"Name":    "Edge",
				"Version": version,
				"Gateway": gin.H{
					"Block-Height": info.Height,
					"URL":          srv.client.Gateway,
					"Status":       "ok",
				},
			},
		)
	}
}
