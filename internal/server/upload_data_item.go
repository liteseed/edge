package server

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database/schema"
)

type UploadDataItemResponse struct {
	Id string `json:"id"`
}

// POST /data
func (s *Context) UploadDataItem(c *gin.Context) {
	header, err := verifyHeaders(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	rawData, err := decodeBody(c, *header.ContentLength)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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

	dataItemData, err := base64.RawURLEncoding.DecodeString(dataItem.Data)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("failed to decode data-item data"))
		return
	}
	checksum := calculateChecksum(dataItemData)

	o := &schema.Order{
		ID:       uuid.New(),
		Status:   schema.Queued,
		StoreID:  storeId,
		Checksum: checksum,
	}

	// SAVE TO DATABASE TO TRACK STATUS
	err = s.database.CreateOrder(o)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, UploadDataResponse{Id: o.ID.String()})
}
