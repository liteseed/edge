package bungo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/schema"
)

func errorResponse(c *gin.Context, err string) {
	// client error
	c.JSON(http.StatusBadRequest, schema.ErrorResponse{
		Err: err,
	})
}

func notFoundResponse(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, schema.ErrorResponse{
		Err: err,
	})
}

func internalErrorResponse(c *gin.Context, err string) {
	// internal error
	c.JSON(http.StatusInternalServerError, schema.ErrorResponse{
		Err: err,
	})
}
