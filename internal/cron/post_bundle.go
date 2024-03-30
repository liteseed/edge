package cron

import (
	"log"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/liteseed/edge/internal/database/schema"
)

func parseDataItemFromOrder(c *Context, o *schema.Order) (*types.BundleItem, error) {
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

func (c *Context) postBundle() {
	o, err := c.database.GetOrdersByStatus(schema.Queued)
	if err != nil {
		return
	}

	if len(*o) == 0 {
		log.Println("no dataitem to post")
		return
	}

	dataItems := []types.BundleItem{}

	for _, order := range *o {
		dataItem, err := parseDataItemFromOrder(c, &order)
		if err != nil {
			log.Println(err)
			log.Println("failed to decode:", order.ID)
			continue
		}
		dataItems = append(dataItems, *dataItem)
	}

	bundle, err := utils.NewBundle(dataItems...)
	if err != nil {
		log.Println("failed to bundle:", err)
		return
	}

	tx, err := c.wallet.SendData([]byte(bundle.BundleBinary), []types.Tag{{Name: "Bundle-Format", Value: "binary"}, {Name: "Bundle-Version", Value: "2.0.0"}})
	if err != nil {
		log.Println("failed to upload:", err)
		return
	}

	for _, order := range *o {
		err = c.database.UpdateTransactionID(order.ID, tx.ID)
		if err != nil {
			log.Println(err)
		}
	}
}
