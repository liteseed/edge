package cron

import "github.com/liteseed/edge/internal/database/schema"

// Notify the AO contract of Successful Data Post
func (c *Cron) Notify() {

	o, err := c.database.GetOrders(&schema.Order{Status: schema.Permanent})
	if err != nil {
		c.logger.Error(
			"failed to fetch order from database",
			"error", err,
		)
		return
	}

	updatedOrders := []schema.Order{}
	for _, order := range *o {
		err := c.contract.Notify(order.ID, order.TransactionId)
		if err != nil {
			c.logger.Error(
				"failed to notify",
				"error", err,
				"order_id", order.ID,
				"order_transaction_id", order.TransactionId,
			)
		}
		updatedOrders = append(updatedOrders, schema.Order{ID: order.ID, Status: schema.Reward})
	}
	err = c.database.UpdateOrders(&updatedOrders)
	if err != nil {
		c.logger.Error(
			"failed to update database",
			"error", err,
		)
	}
}
