package cron

import (
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/liteseed/edge/internal/database/schema"
)

func parseDataItemFromOrder(c *Cron, o *schema.Order) (*types.BundleItem, error) {
	rawDataItem, err := c.store.Get(o.ID)
	if err != nil {
		return nil, err
	}
	dataItem, err := utils.DecodeBundleItem(rawDataItem)
	if err != nil {
		return nil, err
	}
	err = c.database.UpdateStatus(o.ID, schema.Sent)
	if err != nil {
		return nil, err
	}
	return dataItem, nil
}

func (c *Cron) postBundle() {
	o, err := c.database.GetOrdersByStatus(schema.Queued)
	if err != nil {
		c.logger.Error(
			"failed to fetch queued orders",
			"error", err,
		)
		return
	}

	if len(*o) == 0 {
		c.logger.Info("no data item to post")
		return
	}

	dataItems := []types.BundleItem{}

	for _, order := range *o {
		dataItem, err := parseDataItemFromOrder(c, &order)
		if err != nil {
			c.logger.Error(
				"failed to decode data item",
				"error", err,
				"order", order.ID,
			)
			continue
		}
		dataItems = append(dataItems, *dataItem)
	}

	bundle, err := utils.NewBundle(dataItems...)
	if err != nil {
		c.logger.Error(
			"failed to bundle data items",
			"error", err,
		)
		return
	}

	tx, err := c.wallet.SendData([]byte(bundle.BundleBinary), []types.Tag{{Name: "Bundle-Format", Value: "binary"}, {Name: "Bundle-Version", Value: "2.0.0"}, {Name: "App-Name", Value: "Edge"}})
	if err != nil {
		c.logger.Error(
			"failed to upload bundle",
			"error", err,
		)
		return
	}

	updatedOrders := []schema.Order{}
	for _, order := range *o {
		updatedOrders = append(updatedOrders, schema.Order{ID: order.ID, Status: schema.Sent, TransactionId: tx.ID})
	}
	err = c.database.UpdateOrder(&updatedOrders)
	if err != nil {
		c.logger.Error(
			"failed to update database",
			"error", err,
		)
	}
}
