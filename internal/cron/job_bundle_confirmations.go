package cron

import (
	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Cron) checkBundleConfirmations(bundleID string) *schema.Order {
	status, err := c.wallet.Client.GetTransactionStatus(bundleID)
	if err != nil {
		c.logger.Error("fail: gateway - get transaction status", "err", err)
		return nil
	}
	if status.NumberOfConfirmations >= 25 {
		return &schema.Order{Status: schema.Confirmed}
	}
	return nil
}

// Check status of the upload on Arweave
func (crn *Cron) JobBundleConfirmations() {
	orders, err := crn.database.GetOrders(&schema.Order{Status: schema.Sent})
	if err != nil {
		crn.logger.Error("fail: database - get orders", "error", err)
		return
	}

	for _, order := range *orders {
		u := crn.checkBundleConfirmations(order.BundleID)
		if u != nil {
			err = crn.database.UpdateOrder(order.ID, u)
			if err != nil {
				crn.logger.Error("fail: database - update order", "err", err)
			}
		}
		continue
	}
}
