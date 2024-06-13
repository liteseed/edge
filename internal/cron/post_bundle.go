package cron

import (
	"github.com/liteseed/edge/internal/database/schema"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction"
	"github.com/liteseed/goar/transaction/bundle"
	"github.com/liteseed/goar/transaction/data_item"
)

func (crn *Cron) parseDataItemFromOrder(o *schema.Order) (*data_item.DataItem, error) {
	raw, err := crn.store.Get(o.ID)
	if err != nil {
		return nil, err
	}

	dataItem, err := data_item.Decode(raw)
	if err != nil {
		return nil, err
	}

	err = crn.database.UpdateOrder(&schema.Order{ID: o.ID, Status: schema.Posted})
	if err != nil {
		return nil, err
	}
	return dataItem, nil
}

func (crn *Cron) PostBundle() {
	orders, err := crn.database.GetOrders(&schema.Order{Payment: schema.Paid})
	if err != nil {
		crn.logger.Error(
			"failed to fetch queued orders",
			"error", err,
		)
		return
	}

	if len(*orders) == 0 {
		crn.logger.Info("no data item to post")
		return
	}

	dataItems := []data_item.DataItem{}

	for _, order := range *orders {
		dataItem, err := crn.parseDataItemFromOrder(&order)
		if err != nil {
			crn.logger.Error(
				"failed to decode data item",
				"error", err,
				"order", order.ID,
			)
			continue
		}
		dataItems = append(dataItems, *dataItem)
	}

	bundle, err := bundle.New(&dataItems)

	if err != nil {
		crn.logger.Error(
			"failed to bundle data items",
			"error", err,
		)
		return
	}
	tx := &transaction.Transaction{
		Data: bundle.RawData,
		Tags: []tag.Tag{{Name: "Bundle-Format", Value: "binary"}, {Name: "Bundle-Version", Value: "2.0.0"}, {Name: "App-Name", Value: "Edge"}},
	}
	err = tx.Sign(crn.signer)
	if err != nil {
		crn.logger.Error(
			"failed to sign transaction",
			"error", err,
		)
		return
	}

	code, err := crn.client.SubmitTransaction(tx)
	if err != nil {
		crn.logger.Error(
			"failed to sign transaction",
			"error", err,
			"code", code,
		)
		return
	}

	for _, order := range *orders {
		err = crn.database.UpdateOrder(&schema.Order{ID: order.ID, Status: schema.Posted, BundleID: tx.ID})
		if err != nil {
			crn.logger.Error(
				"failed to update database",
				"error", err,
			)
		}
	}

}
