package cron

import (
	"log"

	"github.com/liteseed/edge/internal/database/schema"
)

// Check status of the upload on Arweave

func (c *Context) SyncStatus() {
	o, err := c.database.GetOrdersByStatus(schema.Sent)
	if err != nil {
		log.Println(err)
		return
	}

	for _, order := range *o {
		status, err := c.wallet.Client.GetTransactionStatus(order.TransactionID)
		if err != nil {
			log.Println(err)
			continue
		}
		if status.NumberOfConfirmations > 10 {
			err = c.database.UpdateStatus(order.ID, schema.Permanent)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
