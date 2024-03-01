package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liteseed/bungo/internal/database/schema"
)

type GetDataResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

// GetData - Get status of data sent to upload
//
// GET /:id
func (api *Routes) GetData(c *gin.Context) {
	id := c.Param("id")
	o, err := api.database.GetOrder(id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, &GetDataResponse{
		Id:     id,
		Status: string(o.Status),
	})
}

type PostDataResponse struct {
	Id string `json:"id"`
}

// PostData
//
// POST /data
func (api *Routes) PostData(c *gin.Context) {
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

		id, err := api.store.Put(data)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		api.database.CreateStore(&schema.Store{
			ID:      id,
			OrderID: o.ID,
		})

	}
	c.JSON(http.StatusOK, PostDataResponse{Id: o.ID.String()})
}
