package cron

import "github.com/liteseed/edge/internal/database/schema"

func (c *Cron) ReleaseReward() {
	o, err := c.database.GetOrders(&schema.Order{Status: schema.Reward})
	if err != nil {
		c.logger.Error(
			"failed to fetch reward orders",
			"error", err,
		)
		return
	}

	updatedOrders := []schema.Order{}
	for _, order := range *o {
		err = c.contract.Release(order.ID)
		if err != nil {
			c.logger.Error(
				"failed to release reward",
				"error", err,
			)
		}
		updatedOrders = append(updatedOrders, schema.Order{ID: order.ID, Status: schema.Done})
	}

	err = c.database.UpdateOrders(&updatedOrders)
	if err != nil {
		c.logger.Error(
			"failed to update database",
			"error", err,
		)
	}
}
