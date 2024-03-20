package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database/schema"
)

type UploadDataItemRequestHeader struct {
	ContentType   *string `header:"content-type" binding:"required"`
	ContentLength *int    `header:"content-length" binding:"required"`
}

type UploadDataItemResponse struct {
	Id string `json:"id"`
}

// POST /data
func (s *Context) UploadDataItem(c *gin.Context) {
	header := &UploadDataRequestHeader{}
	if err := c.ShouldBindHeader(header); err != nil {
		c.JSON(400, err.Error())
		return
	}
	if *header.ContentLength == 0 || *header.ContentLength > MAX_DATA_SIZE {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("content-length: supported range 1B - %dB", MAX_DATA_SIZE))
		return
	}
	if *header.ContentType != CONTENT_TYPE_OCTET_STREAM {
		c.AbortWithError(http.StatusBadRequest, errors.New("content-type: unsupported"))
		return
	}

	rawData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("body: failed to parse"))
		return
	}
	if len(rawData) == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("body: required"))
		return
	}
	if len(rawData) != *header.ContentLength {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("content-length, body: length mismatch (%d, %d)", *header.ContentLength, len(rawData)))
		return
	}

	dataItem, err := transaction.DecodeDataItem(rawData)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid data item"))
		return
	}

	valid, err := transaction.VerifyDataItem(dataItem)
	if !valid || err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid data item"))
		return
	}

	storeId, err := s.store.Put(rawData)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("failed to save data"))
		return
	}

	o := &schema.Order{
		ID:      uuid.New(),
		Status:  schema.Queued,
		StoreID: storeId,
	}

	// SAVE TO DATABASE TO TRACK STATUS
	err = s.database.CreateOrder(o)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, UploadDataResponse{Id: o.ID.String()})
}
