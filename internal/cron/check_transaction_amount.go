package cron

import (
	"strconv"

	"github.com/liteseed/edge/internal/database/schema"
)

// Check Transaction ID
// Price of Upload
// Number of Confirmation > 10
func (crn *Cron) CheckTransactionAmount() {
	orders, err := crn.database.GetOrders(&schema.Order{Payment: schema.Confirmed})
	if err != nil {
		crn.logger.Error("fail: database - get orders", "error", err)
		return
	}
	for _, order := range *orders {
		o := schema.Order{ID: order.ID}

		tx, err := crn.client.GetTransactionByID(order.TransactionID)
		if err != nil {
			crn.logger.Error("fail: gateway - get transaction by id", "err", err)
			continue
		}

		payment, err := strconv.ParseUint(tx.Quantity, 10, 64)
		if err != nil {
			crn.logger.Error("fail: internal - conversion to uint", "err", err)
			continue
		}

		res, err := crn.client.GetTransactionPrice(order.Size, "")
		if err != nil {
			crn.logger.Error("fail: gateway - get transaction status", "err", err)
			continue
		}

		price, err := strconv.ParseUint(res, 10, 64)
		if err != nil {
			crn.logger.Error("fail: internal - conversion to uint", "err", err)
			continue
		}

		if payment >= price && tx.Target == crn.signer.Address {
			o.Payment = schema.Paid
		} else {
			o.Payment = schema.Invalid
			o.Status = schema.Failed
		}

		err = crn.database.UpdateOrder(&o)
		if err != nil {
			return
		}
	}

}
