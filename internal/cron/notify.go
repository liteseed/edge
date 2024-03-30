package cron

import (
	"log"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/liteseed/edge/internal/database/schema"
)

// Notify the AO contract of Successful Data Post

const PROCESS = "lJLnoDsq8z0NJrTbQqFQ1arJayfuqWPqwRaW_3aNCgk"

func (c *Context) notify() {

	o, err := c.database.GetOrdersByStatus(schema.Permanent)
	if err != nil {
		log.Println(err)
		return
	}

	itemSigner, err := goar.NewItemSigner(c.wallet.Signer)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	for _, order := range *o {
		message, err := c.ao.SendMessage(PROCESS, "notify", []types.Tag{{Name: "Action", Value: "Notify"}, {Name: "Transaction", Value: order.ID}, {Name: "Status", Value: "1"}}, "", itemSigner)
		if err != nil {
			log.Println(err, message)
			continue
		}
		err = c.database.UpdateStatus(order.ID, schema.Reward)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
