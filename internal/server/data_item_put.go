package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/edge/internal/database/schema"
)

type DataItemPutResponse struct {
	ID        string `json:"id"`
	PaymentID string `json:"payment_id"`
}

// PUT /tx/:id/:payment_id
func (srv *Server) DataItemPut(ctx *gin.Context) {
	ID := ctx.Param("id")
	payment := ctx.Param("payment_id")

	err := srv.contract.Pay(ID, payment)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create order"})
		return
	}

	err = srv.database.UpdateOrder(ID, &schema.Order{TransactionID: payment, Status: schema.Queued})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create order"})
		return
	}

	ctx.JSON(http.StatusAccepted, &DataItemPutResponse{ID: ID, PaymentID: payment})
}
