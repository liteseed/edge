package cron

import (
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
)

// Check status of the upload on Arweave
func (c *Cron) CheckBundleConfirmation() {
	orders, err := c.database.GetOrders(&schema.Order{Status: schema.Release}, database.ConfirmationsLessThan25)
	if err != nil {
		c.logger.Error("fail: database - get orders", "error", err)
		return
	}

	for _, order := range *orders {
		status, err := c.wallet.Client.GetTransactionStatus(order.TransactionID)
		if err != nil {
			c.logger.Error("fail: gateway - get transaction status", "error", err)
			continue
		}
		if status.NumberOfConfirmations >= 25 {
			err = c.database.UpdateOrder(&schema.Order{ID: order.ID, Confirmations: uint(status.NumberOfConfirmations)})
			if err != nil {
				c.logger.Error("fail: database - update order", "err", err)
				continue
			}
		}
	}
}
