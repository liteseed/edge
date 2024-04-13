package cron

import "github.com/liteseed/edge/internal/database/schema"

func (c *Config) ReleaseReward() {
	o, err := c.database.GetOrdersByStatus(schema.Reward)
	if err != nil {
		c.logger.Error(
			"failed to fetch reward orders",
			"error", err,
		)
		return
	}

	for _, order := range *o {
		err = c.contract.Release(order.ID)
		if err != nil {
			c.logger.Error(
				"failed to release reward",
				"error", err,
			)
		}
		err = c.database.UpdateOrder(order.ID, &schema.Order{Status: schema.Done})
		if err != nil {
			c.logger.Error(
				"failed to update order in database",
				"error", err,
			)
		}

	}
}
