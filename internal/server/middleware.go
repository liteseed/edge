package server

import (
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func JSONLogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Println(err)
			}
			c.JSON(-1, c.Errors)
		}
	}

}
