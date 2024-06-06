package cron

import (
	"errors"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Cron) PriceOfUpload(b string, target string) (uint64, error) {
	res, err := http.Get(c.gateway + "/price/" + b + "/" + target)
	if err != nil {
		return 0, err
	}

	r, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	if res.StatusCode >= 400 {
		return 0, errors.New(string(r))
	}

	cost := big.NewInt(0)
	cost.SetString(string(r), 10)

	return cost.Uint64(), nil
}

// Check Transaction ID
// Price of Upload
// Number of Confirmation > 10
func (c *Cron) CheckTransactionAmount() {
	orders, err := c.database.GetOrders(&schema.Order{Payment: schema.Confirmed})
	if err != nil {
		c.logger.Error("fail: database - get orders", "error", err)
		return
	}
	for _, order := range *orders {
		o := schema.Order{ID: order.ID}

		transaction, err := c.wallet.Client.GetTransactionByID(order.TransactionID)
		if err != nil {
			c.logger.Error("fail: gateway - get transaction by id", "err", err)
			continue
		}

		payment, err := strconv.ParseUint(transaction.Quantity, 10, 64)
		if err != nil {
			c.logger.Error("fail: internal - conversion to uint", "err", err)
			continue
		}

		price, err := c.PriceOfUpload(strconv.FormatUint(uint64(order.Size), 10), "")
		if err != nil {
			c.logger.Error("fail: gateway - get transaction status", "err", err)
			continue
		}

		if payment >= price && transaction.Target == c.wallet.Signer.Address {
			o.Payment = schema.Paid
		} else {
			o.Payment = schema.Invalid
			o.Status = schema.Failed
		}
		err = c.database.UpdateOrder(&o)
		if err != nil {
			return
		}
	}

}
