package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"strconv"

	"github.com/everFinance/goar/types"
	"github.com/gin-gonic/gin"
	"github.com/liteseed/edge/internal/contracts"
)

const PROCESS = "lJLnoDsq8z0NJrTbQqFQ1arJayfuqWPqwRaW_3aNCgk"

type UploadRequestHeader struct {
	ContentType   *string `header:"content-type" binding:"required"`
	ContentLength *int    `header:"content-length" binding:"required"`
}

func verifyHeaders(c *gin.Context) (*UploadRequestHeader, error) {
	header := &UploadRequestHeader{}
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

func decodeBody(c *gin.Context, contentLength int) ([]byte, error) {
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

func verify(contract *contracts.Context, dataItem *types.BundleItem) (bool, error) {
	res, err := contract.GetUpload(dataItem.Id)
	if err != nil {
		return false, err
	}
	if res == nil {
		return false, errors.New("no upload")
	}
	rawData, err := base64.RawURLEncoding.DecodeString(dataItem.Data)
	if err != nil {
		return false, err
	}

	price, err := getPrice(len(rawData))
	if err != nil {
		return false, err
	}

	log.Println(price)
	// if quantity >= 0 {
	// 	return false, errors.New("quantity isn't enough")
	// }

	return true, nil
}

func getPrice(b int) (int, error) {

	res, err := http.Get(fmt.Sprint("https://arweave.net/price/", b))
	if err != nil {
		return 0, err
	}

	r, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil
	}

	i, err := strconv.Atoi(string(r))
	if err != nil {
		return 0, nil
	}

	return i, nil
}
