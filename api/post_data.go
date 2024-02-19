package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type PostDataResponse struct {
	Id   string `json:"id"`
	Size int64  `json:"size"`
}

// PostData
//
// POST /data
func (a *API) PostData(c *gin.Context) {
	log.Println("POST DATA")
	contentLength, err := strconv.Atoi(c.Request.Header.Get("content-length"))
	if err != nil {
		log.Println("request has no content length header!")
	}

	if contentLength > MAX_DATA_ITEM_SIZE {
		BadRequest(
			c,
			fmt.Sprintf("Data item size is currently limited to %d bytes!", MAX_DATA_ITEM_SIZE),
		)
	}
	data, err := io.ReadAll(c.Request.Body)

	id, err := a.store.Save(data)
	c.Render(http.StatusOK, render.JSON{Data: (PostDataResponse{Id: id, Size: int64(contentLength)})})
}
