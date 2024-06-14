package cron

import (
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Cron) checkBundleConfirmations(ID string, transactionID string) *schema.Order {
	status, err := c.wallet.Client.GetTransactionStatus(transactionID)
	if err != nil {
		c.logger.Error("fail: gateway - get transaction status", "err", err)
		return nil
	}
	if status.NumberOfConfirmations >= 10 {
		return &schema.Order{
			ID:     ID,
			Status: schema.Queued,
		}
	}
	return nil
}

// Check status of the upload on Arweave
func (crn *Cron) JobBundleConfirmations() {
	orders, err := crn.database.GetOrders(&schema.Order{Status: schema.Release}, database.ConfirmationsLessThan25)
	if err != nil {
		crn.logger.Error("fail: database - get orders", "error", err)
		return
	}

	for _, order := range *orders {
		u := crn.checkBundleConfirmations(order.ID, order.TransactionID)
		if u != nil {
			err = crn.database.UpdateOrder(u)
			if err != nil {
				crn.logger.Error("fail: database - update order", "err", err)
			}
		}
		continue
	}
}
