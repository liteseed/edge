package routes

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database/schema"
)

type UploadDataResponse struct {
	Id string `json:"id"`
}

// POST /data
func (api *Routes) UploadData(c *gin.Context) {
	contentLength, err := strconv.Atoi(c.Request.Header.Get("content-length"))
	if err != nil {
		log.Println("request has no content length header!")
	}

	if contentLength > MAX_DATA_SIZE {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	contentType := c.Request.Header.Get("content-type")
	if contentType == "" {
		log.Println("request has no content type")
	} else if contentType != CONTENT_TYPE_OCTET_STREAM {
		c.AbortWithError(http.StatusBadRequest, c.Error(errors.New("unexpected content type")))
		return
	}

	rawData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, c.Error(errors.New("unable to decode data")))
		return
	}

	dataItem, err := transaction.NewDataItem(rawData, *api.signer, "", "", []transaction.Tag{})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, c.Error(errors.New("could not create a data item")))
		return
	}

	storeId, err := api.store.Put(dataItem.Raw)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, c.Error(errors.New("failed to save data")))
		return
	}

	// ADD TO NEXT BUNDLE
	o := &schema.Order{
		ID:      uuid.New(),
		Status:  schema.Queued,
		StoreID: storeId,
	}

	// SAVE TO DATABASE TO TRACK STATUS
	err = api.database.CreateOrder(o)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, UploadDataResponse{Id: o.ID.String()})
}
