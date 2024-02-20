package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDataResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

// GetData - Get status of data sent to upload
//
// GET /:id
func (api *API) GetData(c *gin.Context) {
	id := c.Param("id")
	o, err := api.database.GetOrder(id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, &GetDataResponse{
		Id:     id,
		Status: string(o.Status),
	})
}
