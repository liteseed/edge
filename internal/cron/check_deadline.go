package cron

import (
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Cron) CheckDeadline() {
	info, err := c.wallet.Client.GetInfo()
	if err != nil {
		c.logger.Error(
			"failed to query gateway",
			"error", err,
		)
		return
	}
	orders, err := c.database.GetOrders(&schema.Order{Status: schema.Created}, database.DeadlinePassed(info.Height))
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
		err := c.database.DeleteOrder(order.ID)
		if err != nil {
			c.logger.Info("unable to remove order")
			continue
		}
	}

}
