package server

import (
	"net/http"

	"github.com/everFinance/goar/utils"
	"github.com/gin-gonic/gin"

	"github.com/liteseed/edge/internal/database/schema"
)

// POST /data-item
func (s *Config) uploadDataItem(c *gin.Context) {
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

	valid, err := verify(s.contract, dataItem)
	if !valid || err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = s.store.Set(dataItem.Id, dataItem.ItemBinary)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	o := &schema.Order{
		ID:     dataItem.Id,
		Status: schema.Queued,
	}

	err = s.database.CreateOrder(o)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
