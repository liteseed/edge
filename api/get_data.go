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
func (a *API) GetData(c *gin.Context) {
	id := c.Param("id")
	o, err := a.db.GetOrder(id)
	if err != nil {
		NotFound(c)
	}
	res := &GetDataResponse{
		Id:     id,
		Status: string(o.Status),
	}
	c.JSON(http.StatusOK, res)
}
