package server

import (
	"net/http"

	"github.com/everFinance/goar/utils"
	"github.com/gin-gonic/gin"

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
	dataItem, err := utils.DecodeBundleItem(rawData)
	if err != nil {
		s.logger.Error(
			"failed to decode data item",
			"error", err,
		)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// err = utils.VerifyBundleItem(*dataItem)
	// if err != nil {
	// 	log.Println("data-item: failed to verify", err)
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, err)
	// 	return
	// }

	valid, err := checkUploadOnContract(s.contract, dataItem)
	if !valid || err != nil {
		s.logger.Error(
			"failed to fetch verify on AO",
			"error", err,
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = s.store.Put(dataItem.Id, dataItem.ItemBinary)
	if err != nil {
		s.logger.Error(
			"failed to save to store",
			"error", err,
		)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	o := &schema.Order{
		ID:     dataItem.Id,
		Status: schema.Queued,
	}

	err = s.database.CreateOrder(o)
	if err != nil {
		s.logger.Error(
			"failed to create order",
			"error", err,
		)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
