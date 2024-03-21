package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database/schema"
)

type UploadDataResponse struct {
	Id string `json:"id"`
}

// POST /data
func (s *Context) UploadData(c *gin.Context) {
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

	dataItem, err := transaction.NewDataItem(rawData, signer.Signer{S: s.signer}, "", "", []transaction.Tag{})
	if err != nil {
		log.Println("data-item: failed to create", err)
		c.AbortWithStatus(http.StatusInternalServerError)
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
		Checksum: checksum,
	}

	// SAVE TO DATABASE TO TRACK STATUS
	err = s.database.CreateOrder(o)
	if err != nil {
		log.Println("database: create order failed", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, UploadDataResponse{Id: o.ID.String()})
}
