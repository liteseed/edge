package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/everFinance/goar/utils"
	"github.com/gin-gonic/gin"

	"github.com/liteseed/edge/internal/database/schema"
)

type DataItemPostRequestHeader struct {
	ContentType   *string `header:"content-type" binding:"required"`
	ContentLength *int    `header:"content-length" binding:"required"`
}

type DataItemPostResponse struct {
	ID                  string   `json:"id"`
	Owner               string   `json:"owner"`
	DataCaches          []string `json:"dataCaches"`
	DeadlineHeight      uint     `json:"deadlineHeight"`
	FastFinalityIndexes []string `json:"fastFinalityIndexes"`
	Price               uint64   `json:"price"`
	Version             string   `json:"version"`
}

func parseHeaders(c *gin.Context) (*DataItemPostRequestHeader, error) {
	header := &DataItemPostRequestHeader{}
	if err := c.ShouldBindHeader(header); err != nil {
		return nil, err
	}
	if *header.ContentLength == 0 || *header.ContentLength > MAX_DATA_ITEM_SIZE {
		return nil, fmt.Errorf("content-length: supported range 1B - %dB", MAX_DATA_ITEM_SIZE)
	}
	if *header.ContentType != CONTENT_TYPE_OCTET_STREAM {
		return nil, fmt.Errorf("content-type: unsupported")
	}
	return header, nil
}

func parseBody(c *gin.Context, contentLength int) ([]byte, error) {
	rawData, err := c.GetRawData()
	if err != nil {
		return nil, err
	}
	if len(rawData) == 0 {
		return nil, errors.New("body: required")
	}
	if len(rawData) != contentLength {
		return nil, fmt.Errorf("content-length, body: length mismatch (%d, %d)", contentLength, len(rawData))
	}
	return rawData, nil
}

// POST /data-item
func (s *Server) DataItemPost(c *gin.Context) {
	header, err := parseHeaders(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	rawData, err := parseBody(c, *header.ContentLength)
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
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	p, err := s.wallet.Client.GetTransactionPrice(*header.ContentLength, nil)
	if err != nil {
		s.logger.Error(
			"failed to calculate transaction price",
			"error", err,
		)
		c.AbortWithError(http.StatusFailedDependency, err)
		return
	}
	err = utils.VerifyBundleItem(*dataItem)
	if err != nil {
		s.logger.Error(
			"failed to verify data item",
			"error", err,
		)
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
		Price:  uint64(p),
	}

	err = s.database.CreateOrder(o)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		&DataItemPostResponse{
			ID:      o.ID,
			Owner:   s.wallet.Signer.Address,
			Price:   o.Price,
			Version: "1.0.0",
		},
	)
}
