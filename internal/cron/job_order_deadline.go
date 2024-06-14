package cron

import (
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
)

func (crn *Cron) JobOrderDeadline() {
	info, err := crn.wallet.Client.GetNetworkInfo()
	if err != nil {
		crn.logger.Error("failed to query gateway", "error", err)
		return
	}
	orders, err := crn.database.GetOrders(&schema.Order{Status: schema.Created}, database.DeadlinePassed(info.Height))
	if err != nil {
		crn.logger.Error("failed to fetch queued orders", "error", err)
		return
	}

	if len(*orders) == 0 {
		crn.logger.Info("no data item to delete")
		return
	}

	for _, order := range *orders {
		err := crn.database.DeleteOrder(order.ID)
		if err != nil {
			crn.logger.Error("unable to remove order", "err", err)
			continue
		}
	}

}
