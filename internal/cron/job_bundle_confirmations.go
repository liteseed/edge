package cron

import "github.com/liteseed/edge/internal/database/schema"


// Check status of the upload on Arweave
func (crn *Cron) JobBundleConfirmations() {
	orders, err := crn.database.GetOrders(&schema.Order{Status: schema.Sent})
	if err != nil {
		crn.logger.Error("fail: database - get orders", "error", err)
		return
	}

	for _, order := range *orders {
		status, err := crn.wallet.Client.GetTransactionStatus(order.BundleID)
		if err != nil {
			crn.logger.Error("fail: gateway - get transaction status", "err", err)
			continue
		}
		if status.NumberOfConfirmations >= 25 {
			err = crn.database.UpdateOrder(order.ID, &schema.Order{Status: schema.Confirmed})
			if err != nil {
				crn.logger.Error("fail: database - update order", "err", err)
				continue
			}
		}
	}
}
