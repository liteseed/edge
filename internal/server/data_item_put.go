package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/liteseed/edge/internal/database/schema"
)

type DataItemPutResponse struct {
	ID             string `json:"id"`
	DeadlineHeight uint   `json:"deadlineHeight"`
}

// PUT /tx/:id/:transaction_id
func (srv *Server) DataItemPut(context *gin.Context) {
	ID := context.Param("id")
	transactionID := context.Param("transaction_id")
	info, err := srv.wallet.Client.GetNetworkInfo()
	if err != nil {
		context.Status(http.StatusFailedDependency)
		log.Println(err)
		return
	}

	deadline := uint(info.Height) + 200
	err = srv.database.UpdateOrder(&schema.Order{ID: ID, TransactionID: transactionID, Status: schema.Queued, DeadlineHeight: deadline})
	if err != nil {
		context.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	context.JSON(
		http.StatusAccepted,
		&DataItemPutResponse{
			ID:             ID,
			DeadlineHeight: deadline,
		},
	)
}
