package cron

import (
	"github.com/liteseed/edge/internal/database/schema"
)

// Notify the AO contract of Successful Data Post
func (crn *Cron) JobPostUpdate() {
	o, err := crn.database.GetOrders(&schema.Order{Status: schema.Confirmed})
	if err != nil {
		crn.logger.Error("failed to fetch order from database", "error", err)
		return
	}

	for _, order := range *o {
		err := crn.contract.Posted(order.ID)
		if err != nil {
			crn.logger.Error("failed to post transaction to contract", "error", err, "order_id", order.ID, "order_transaction_id", order.TransactionID)
			continue
		}
		err = crn.database.UpdateOrder(order.ID, &schema.Order{Status: schema.Release})
		if err != nil {
			crn.logger.Error(
				"failed to update database",
				"error", err,
			)
			continue
		}
	}
}
