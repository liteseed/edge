package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) PriceGet(c *gin.Context) {
	b, isSet := c.Params.Get("bytes")
	if isSet || b == "" || b == "0" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} 
	resp, err := http.Get(fmt.Sprintf("https://arweave.net/price/%s", b))
	if err != nil {
		c.AbortWithStatus(http.StatusFailedDependency)
		return
	}
	
	c.JSON(200, resp.Body)
}
