package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database/schema"
)

// POST /data-item
func (s *Context) uploadDataItem(c *gin.Context) {
	id := c.Param("id")
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
		log.Println("data-item: failed to create", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	valid, err := transaction.VerifyDataItem(dataItem)
	if !valid || err != nil {
		log.Println("data-item: failed to verify", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	storeId, err := s.store.Put(dataItem.Raw)
	if err != nil {
		log.Println("store: failed to save", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	checksum := calculateChecksum(rawData)

	// ADD TO NEXT BUNDLE
	o := &schema.Order{
		ID:       uuid.New(),
		Status:   schema.Queued,
		StoreID:  storeId,
		PublicID: id,
		Checksum: checksum,
	}

	// SAVE TO DATABASE TO TRACK STATUS
	err = s.database.CreateOrder(o)
	if err != nil {
		log.Println("database: create order failed", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
