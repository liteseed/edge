package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, err)
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, errors.New("not found"))
}

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, errors.New("something went wrong"))
}
