package cron

import (
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Cron) JobRemoveBundle() {
	info, err := c.wallet.Client.GetNetworkInfo()
	if err != nil {
		c.logger.Error(
			"failed to query gateway",
			"error", err,
		)
		return
	}

	orders, err := c.database.GetOrders(&schema.Order{Status: schema.Invalid}, database.DeadlinePassed(info.Height))
	if err != nil {
		c.logger.Error(
			"failed to fetch queued orders",
			"error", err,
		)
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
