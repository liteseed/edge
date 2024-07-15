package cron

import "github.com/liteseed/edge/internal/database/schema"

func (crn *Cron) JobDeleteBundle() {
	orders, err := crn.database.GetOrders(&schema.Order{Status: schema.Created})
	if err != nil {
		crn.logger.Error("failed to fetch queued orders", "err", err)
		return
	}

	if len(*orders) == 0 {
		crn.logger.Info("no data item to delete")
		return
	}

	for _, order := range *orders {
		n, err := crn.wallet.Client.GetNetworkInfo()
		if err != nil {
			crn.logger.Error("failed to fetch network info", "err", err)
			continue
		}
		if n.Height >= int64(order.DeadlineHeight)+100 {
			err = crn.database.DeleteOrder(order.ID)
			if err != nil {
				crn.logger.Error(
					"failed to delete order from database",
					"error", err,
				)
				continue
			}
		}

	}
}
