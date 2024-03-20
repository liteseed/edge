package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			log.Println(err)
		}
		c.JSON(-1, c.Errors)
	}
}
