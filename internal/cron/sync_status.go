package cron

import "github.com/liteseed/edge/internal/database/schema"

// Check status of the upload on Arweave
func (c *Cron) SyncStatus() {
	orders, err := c.database.GetOrders(&schema.Order{Status: schema.Sent})
	if err != nil {
		c.logger.Error(
			"failed to fetch sent orders",
			"error", err,
		)
		return
	}

	for _, order := range *orders {
		status, err := c.wallet.Client.GetTransactionStatus(order.TransactionId)
		if err != nil {
			c.logger.Error(
				"failed to fetch transaction status",
				"error", err,
			)
			continue
		}
		if status.NumberOfConfirmations > 10 {
			err = c.database.UpdateOrder(&schema.Order{ID: order.ID, Status: schema.Permanent})
			if err != nil {
				c.logger.Error(
					"failed to fetch transaction status",
					"error", err,
				)
				continue
			}
		}
	}
}
