package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func JSONLogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("error", err.Err)
			}
			c.JSON(-1, c.Errors)
		}
	}
}
