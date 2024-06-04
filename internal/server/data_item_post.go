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
	DeadlineHeight      uint   `json:"deadlineHeight"`
	FastFinalityIndexes []string `json:"fastFinalityIndexes"`
	Version             string   `json:"version"`
}

func parseHeaders(context *gin.Context) (*DataItemPostRequestHeader, error) {
	header := &DataItemPostRequestHeader{}
	if err := context.ShouldBindHeader(header); err != nil {
		return nil, err
	}
	if *header.ContentLength == 0 || uint(*header.ContentLength) > MAX_DATA_ITEM_SIZE {
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

// POST /tx
func (s *Server) DataItemPost(context *gin.Context) {
	header, err := parseHeaders(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawData, err := parseBody(context, *header.ContentLength)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataItem, err := utils.DecodeBundleItem(rawData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode bundle"})
		return
	}

	err = utils.VerifyBundleItem(*dataItem)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "failed to verify bundle"})
		return
	}

	owner, err := utils.OwnerToAddress(dataItem.Owner)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, err)
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
		return
	}

	info, err := s.wallet.Client.GetInfo()
	if err != nil {
		context.JSON(http.StatusFailedDependency, gin.H{"error": "failed to query gateway"})
		return
	}
	deadline := uint(info.Height) + 200
	o := &schema.Order{
		ID:             dataItem.Id,
		Status:         schema.Created,
		Price:          uint(p),
		DeadlineHeight: deadline,
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
			Owner:               owner,
			Version:             "1.0.0",
			DeadlineHeight:      deadline,
			DataCaches:          []string{s.gateway},
			FastFinalityIndexes: []string{s.gateway},
		},
	)
}
