package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database/schema"
)

// POST /data-item
func (s *Context) uploadDataItem(c *gin.Context) {
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

	valid, err = checkUploadOnContract(s.ao, &signer.Signer{S: s.signer}, dataItem)
	if !valid || err != nil {
		log.Println("data-item: failed to verify on ao", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = s.store.Put(dataItem.ID, dataItem.Raw)
	if err != nil {
		log.Println("store: failed to save", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	checksum := calculateChecksum(rawData)

	o := &schema.Order{
		ID:       dataItem.ID,
		Status:   schema.Queued,
		Checksum: checksum,
	}

	err = s.database.CreateOrder(o)
	if err != nil {
		log.Println("database: create order failed", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
