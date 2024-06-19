package cron

import (
	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Cron) JobReleaseReward() {
	o, err := c.database.GetOrders(&schema.Order{Status: schema.Confirmed})
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
			continue
		}
		err = c.database.UpdateOrder(order.ID, &schema.Order{Status: schema.Permanent})
		if err != nil {
			c.logger.Error(
				"failed to update database",
				"error", err,
			)
			continue
		}
	}

}
