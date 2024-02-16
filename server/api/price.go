package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/everFinance/goar/types"
	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/schema"
)

// Get the current price to post a transaction onto Arweave
//
// GET price/:size
func GetTransactionPrice(router *gin.RouterGroup) {
	router.GET("/price/:size", func(c *gin.Context) {
		size, err := strconv.ParseInt(c.Param("size"), 10, 64)
		if err != nil {
			errorResponse(c, err.Error())
		}
		// price = chunkNum*deltaPrice(fee for per chunk) + basePrice
		price := calculatePrice(schema.ArFee{Base: 100, PerChunk: 100}, size)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("%d", price)))
	})
}

func calculatePrice(fee schema.ArFee, dataSize int64) int64 {
	count := int64(0)
	if dataSize > 0 {
		count = (dataSize-1)/types.MAX_CHUNK_SIZE + 1
	}
	totPrice := fee.Base + count*fee.PerChunk
	return totPrice
}
