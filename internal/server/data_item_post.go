package server

import (
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

func parseHeaders(context *gin.Context) (*DataItemPostRequestHeader, error) {
	header := &DataItemPostRequestHeader{}
	if err := context.ShouldBindHeader(header); err != nil {
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

func parseBody(context *gin.Context, contentLength int) ([]byte, error) {
	rawData, err := context.GetRawData()
	if err != nil {
		return nil, err
	}
	if len(rawData) != contentLength {
		return nil, fmt.Errorf("content-length, body: length mismatch (%d, %d)", contentLength, len(rawData))
	}
	return rawData, nil
}

// POST /data-item
func (s *Server) DataItemPost(context *gin.Context) {
	header, err := parseHeaders(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Error(err)
		return
	}

	rawData, err := parseBody(context, *header.ContentLength)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Error(err)
		return
	}

	dataItem, err := utils.DecodeBundleItem(rawData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode bundle"})
		context.Error(err)
		return
	}

	err = utils.VerifyBundleItem(*dataItem)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "failed to verify bundle"})
		context.Error(err)
		return
	}

	err = s.store.Set(dataItem.Id, dataItem.ItemBinary)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	
	p, err := s.wallet.Client.GetTransactionPrice(*header.ContentLength, nil)
	if err != nil {
		context.JSON(http.StatusFailedDependency, gin.H{"error": "failed to query gateway"})
		context.Error(err)
		return
	}

	info, err := s.wallet.Client.GetInfo()
	if err != nil {
		context.JSON(http.StatusFailedDependency, gin.H{"error": "failed to query gateway"})
		context.Error(err)
		return
	}

	o := &schema.Order{
		ID:     dataItem.Id,
		Status: schema.Queued,
		Price:  uint64(p),
	}

	err = s.database.CreateOrder(o)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	context.JSON(
		http.StatusCreated,
		&DataItemPostResponse{
			ID:                  o.ID,
			Owner:               s.wallet.Signer.Address,
			Price:               o.Price,
			Version:             "1.0.0",
			DeadlineHeight:      uint(info.Height + 200),
			DataCaches:          []string{s.gateway},
			FastFinalityIndexes: []string{s.gateway},
		},
	)
}
