package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) PriceGet(c *gin.Context) {
	b, isSet := c.Params.Get("bytes")
	if isSet || b == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	size, err := strconv.ParseUint(b, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	p, err := s.wallet.Client.GetTransactionPrice(int(size), nil)
	if err != nil || size <= 0 {
		c.AbortWithStatus(http.StatusFailedDependency)
		return
	}
	c.JSON(200, p)
}
