package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/liteseed/bungo/database/schema"
)

type PostDataResponse struct {
	Id string `json:"id"`
}

// PostData
//
// POST /data
func (api *API) PostData(c *gin.Context) {
	contentLength, err := strconv.Atoi(c.Request.Header.Get("content-length"))
	if err != nil {
		log.Println("request has no content length header!")
	}

	if contentLength > MAX_DATA_ITEM_SIZE {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	files := form.File["files"]

	o := &schema.Order{
		ID:     uuid.New(),
		Status: schema.Queued,
	}
	// SAVE TO DATABASE TO TRACK STATUS
	err = api.database.CreateOrder(o)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, file := range files {
		// SAVE FILE TO OBJECT STORE
		data := []byte{}
		f, err := file.Open()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		f.Read(data)

		id, err := api.store.Save(data)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		api.database.CreateStore(&schema.Store{
			ID:      id,
			OrderID: o.ID,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	c.JSON(http.StatusOK, PostDataResponse{Id: o.ID.String()})
}
