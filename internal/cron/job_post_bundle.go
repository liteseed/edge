package cron

import (
	"github.com/liteseed/edge/internal/database/schema"
	"github.com/liteseed/goar/crypto"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction/bundle"
	"github.com/liteseed/goar/transaction/data_item"
	"log"
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
	return dataItem, nil
}

func (crn *Cron) JobPostBundle() {
	orders, err := crn.database.GetOrders(&schema.Order{Status: schema.Queued})
	if err != nil {
		crn.logger.Error("failed to fetch queued orders", "error", err)
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
			crn.logger.Error("failed to decode data item", "error", err, "order", order.ID)
			continue
		}
		dataItems = append(dataItems, *dataItem)
	}

	b, err := bundle.New(&dataItems)
	if err != nil {
		crn.logger.Error("fail - internal - bundle data items", "error", err)
		return
	}

	tx := crn.wallet.CreateTransaction(b.Raw, "", "", &[]tag.Tag{{Name: "Bundle-Format", Value: "binary"}, {Name: "Bundle-Version", Value: "2.0.0"}, {Name: "App-Name", Value: "Edge"}})

	_, err = crn.wallet.SignTransaction(tx)
	if err != nil {
		crn.logger.Error("fail - internal - sign transaction", "err", err)
		return
	}

	data, err := crypto.Base64URLDecode(tx.Data)
	if err != nil {
		crn.logger.Error("failed to decode transaction data", "err", err)
		return
	}

	tx.Data = ""
	res, err := crn.wallet.Client.SubmitTransaction(tx)
	log.Println(res)
	if err != nil {
		crn.logger.Error("failed to send transaction", "err", err)
		return
	}

	for i := 0; i < len(tx.ChunkData.Chunks); i++ {
		c, err := tx.GetChunk(i, data)
		if err != nil {
			crn.logger.Error("failed to upload chunk", "err", err)
			return
		}
		res, err = crn.wallet.Client.UploadChunk(c)
		if err != nil {
			crn.logger.Error("failed to upload chunk", "err", err)
			return
		}
	}

	for _, d := range dataItems {
		err = crn.database.UpdateOrder(d.ID, &schema.Order{Status: schema.Sent, BundleID: tx.ID})
		if err != nil {
			crn.logger.Error("failed to update database", "err", err)
		}
	}
}
