package cron

import "github.com/liteseed/edge/internal/database/schema"

// Check status of the upload on Arweave

func (c *Config) SyncStatus() {
	o, err := c.database.GetOrdersByStatus(schema.Sent)
	if err != nil {
		c.logger.Error(
			"failed to fetch sent orders",
			"error", err,
		)
		return
	}

	for _, order := range *o {
		status, err := c.wallet.Client.GetTransactionStatus(order.TransactionID)
		if err != nil {
			c.logger.Error(
				"failed to fetch transaction status",
				"error", err,
			)
			continue
		}
		if status.NumberOfConfirmations > 10 {
			err = c.database.UpdateStatus(order.ID, schema.Permanent)
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
