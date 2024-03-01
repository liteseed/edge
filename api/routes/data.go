package routes

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/everFinance/goar/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liteseed/bungo/internal/database/schema"
)

const (
	CONTENT_TYPE_OCTET_STREAM = "application/octet-stream"
)

type GetDataResponse struct {
	Id     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

// GetData - Get status of data sent to upload
//
// GET /:id
func (api *Routes) GetData(c *gin.Context) {
	param := c.Param("id")
	id, err := uuid.Parse(param)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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

type PostDataResponse struct {
	Id string `json:"id"`
}

// PostData
//
// POST /data
func (api *Routes) PostData(c *gin.Context) {
	contentLength, err := strconv.Atoi(c.Request.Header.Get("content-length"))
	if err != nil {
		log.Println("request has no content length header!")
	}

	if contentLength > MAX_DATA_ITEM_SIZE {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	contentType := c.Request.Header.Get("content-type")
	if contentType == "" {
		log.Println("request has no content type")
	} else if contentType != CONTENT_TYPE_OCTET_STREAM {
		c.AbortWithError(http.StatusBadRequest, c.Error(errors.New("unexpected content type")))
		return
	}
	bundle := &types.Bundle{}
	api.store.Put(bundle.BundleBinary)

	o := &schema.Order{
		ID:     uuid.New(),
		Status: schema.Queued,
	}
	// SAVE TO DATABASE TO TRACK STATUS
	err = api.database.CreateOrder(o)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, PostDataResponse{Id: o.ID.String()})
}
