package cron

import "github.com/liteseed/edge/internal/database/schema"

// Notify the AO contract of Successful Data Post

const PROCESS = "lJLnoDsq8z0NJrTbQqFQ1arJayfuqWPqwRaW_3aNCgk"

func (c *Context) notify() {

	o, err := c.database.GetOrdersByStatus(schema.Permanent)
	if err != nil {
		c.logger.Error(
			"failed to fetch order from database",
			"error", err,
		)
		return
	}

	for _, order := range *o {
		err := c.contract.Notify(order.ID, order.TransactionID)
		if err != nil {
			c.logger.Error(
				"failed to notify",
				"error", err,
				"order_id", order.ID,
				"order_transaction_id", order.TransactionID,
			)
			continue
		}
		err = c.database.UpdateStatus(order.ID, schema.Reward)
		if err != nil {
			c.logger.Error(
				"failed to update status of order",
				"error", err,
				"order_id", order.ID,
				"order_transaction_id", order.TransactionID,
			)
			continue
		}
	}
}
