package cron

import (
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
)

// Check status of the upload on Arweave
func (crn *Cron) CheckBundleConfirmation() {
	orders, err := crn.database.GetOrders(&schema.Order{Status: schema.Release}, database.ConfirmationsLessThan25)
	if err != nil {
		crn.logger.Error("fail: database - get orders", "error", err)
		return
	}

	for _, order := range *orders {
		status, err := crn.client.GetTransactionStatus(order.TransactionID)
		if err != nil {
			crn.logger.Error("fail: gateway - get transaction status", "error", err)
			continue
		}
		if status.NumberOfConfirmations >= 25 {
			err = crn.database.UpdateOrder(&schema.Order{ID: order.ID, Confirmations: uint(status.NumberOfConfirmations)})
			if err != nil {
				crn.logger.Error("fail: database - update order", "err", err)
				continue
			}
		}
	}
}
