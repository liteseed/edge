package api

import (
	"fmt"
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
func (a *API) PostData(c *gin.Context) {
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
	form, err := c.MultipartForm()
	if err != nil {
		BadRequest(
			c,
			"Unable to parse data",
		)
	}
	files := form.File["files"]

	o := &schema.Order{
		ID:     uuid.New(),
		Status: schema.Queued,
	}
	// SAVE TO DATABASE TO TRACK STATUS
	err = a.db.CreateOrder(o)
	if err != nil {
		log.Println(err)
		InternalServerError(c)
	}

	for _, file := range files {
		// SAVE FILE TO OBJECT STORE
		data := []byte{}
		f, err := file.Open()
		if err != nil {
			BadRequest(c, err.Error())
		}
		f.Read(data)

		id, err := a.store.Save(data)
		if err != nil {
			log.Println(err)
			InternalServerError(c)
		}
		a.db.CreateStore(&schema.Store{
			ID:      id,
			OrderID: o.ID,
		})
		if err != nil {
			log.Println(err)
			InternalServerError(c)
		}
	}

	c.JSON(http.StatusOK, PostDataResponse{Id: o.ID.String()})
}
