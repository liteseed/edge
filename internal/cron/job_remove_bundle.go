package cron

import "github.com/liteseed/edge/internal/database/schema"

func (c *Cron) JobDeleteBundle() {
	orders, err := c.database.GetOrders(&schema.Order{Status: schema.Permanent})
	if err != nil {
		c.logger.Error("failed to fetch queued orders", err)
		return
	}

	if len(*orders) == 0 {
		c.logger.Info("no data item to delete")
		return
	}

	for _, order := range *orders {
		err = c.database.DeleteOrder(order.ID)
		if err != nil {
			c.logger.Error(
				"failed to delete order from database",
				"error", err,
			)
			continue
		}
	}
}
