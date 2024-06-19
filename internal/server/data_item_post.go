package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/transaction/data_item"

	"github.com/liteseed/edge/internal/database/schema"
)

type DataItemPostRequestHeader struct {
	ContentType   *string `header:"content-type" binding:"required"`
	ContentLength *string `header:"content-length" binding:"required"`
}

type DataItemPostResponse struct {
	ID                  string   `json:"id"`
	Owner               string   `json:"owner"`
	DataCaches          []string `json:"dataCaches"`
	DeadlineHeight      uint     `json:"deadlineHeight"`
	FastFinalityIndexes []string `json:"fastFinalityIndexes"`
	Version             string   `json:"version"`
}

func parseHeaders(context *gin.Context) (*DataItemPostRequestHeader, error) {
	header := &DataItemPostRequestHeader{}
	if err := context.ShouldBindHeader(header); err != nil {
		return nil, errors.New("required header(s) - content-type, content-length")
	}
	if *header.ContentType != CONTENT_TYPE_OCTET_STREAM {
		return nil, errors.New("required header(s) - content-type: application/octet-stream")
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
func (srv *Server) DataItemPost(ctx *gin.Context) {
	header, err := parseHeaders(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contentLength, err := strconv.Atoi(*header.ContentLength)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rawData, err := parseBody(ctx, contentLength)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataItem, err := data_item.Decode(rawData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode data item"})
		return
	}

	err = data_item.Verify(dataItem)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to verify data item"})
		return
	}

	owner, err := crypto.GetAddressFromOwner(dataItem.Owner)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	err = srv.store.Set(dataItem.ID, dataItem.Raw)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	info, err := srv.wallet.Client.GetNetworkInfo()
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, gin.H{"error": "failed to query gateway"})
		return
	}
	deadline := uint(info.Height) + 200
	o := &schema.Order{
		ID:             dataItem.ID,
		Status:         schema.Created,
		Size:           len(dataItem.Raw),
		DeadlineHeight: deadline,
	}

	err = srv.database.CreateOrder(o)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		&DataItemPostResponse{
			ID:                  o.ID,
			Owner:               owner,
			Version:             "1.0.0",
			DeadlineHeight:      deadline,
			DataCaches:          []string{srv.wallet.Client.Gateway},
			FastFinalityIndexes: []string{srv.wallet.Client.Gateway},
		},
	)
}
