package server

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/argo/ao"
	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
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
	if *header.ContentLength == 0 || *header.ContentLength > MAX_DATA_SIZE {
		return nil, fmt.Errorf("content-length: supported range 1B - %dB", MAX_DATA_SIZE)
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

func calculateChecksum(rawData []byte) string {
	rawChecksum := md5.Sum(rawData)
	return hex.EncodeToString(rawChecksum[:])
}

func checkUploadOnContract(ao *ao.AO, s *signer.Signer, dataItem *transaction.DataItem) (bool, error) {
	tags := []transaction.Tag{{Name: "Action", Value: "Upload"}, {Name: "Transaction", Value: dataItem.ID}}
	message, err := ao.SendMessage(PROCESS, "initiate", tags, "", s)
	if err != nil {
		return false, err
	}
	result, err := ao.ReadResult(PROCESS, message)
	if err != nil {
		return false, err
	}

	checksum := ""
	if checksum != result.Messages[0]["Checksum"] {
		return false, errors.New("checksum doesn't match")
	}

	quantity := 100
	if quantity != result.Messages[0]["Quantity"] {
		return false, errors.New("quantity isn't enough")
	}

	return true, nil
}
