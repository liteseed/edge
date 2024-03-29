package cron

import (
	"log"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
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

	for _, order := range *o {
		message, err := c.ao.SendMessage(PROCESS, "notify", []transaction.Tag{{Name: "Action", Value: "Notify"}, {Name: "Transaction", Value: order.ID}, {Name: "Status", Value: "1"}}, "", &signer.Signer{S: c.wallet.Signer})
		if err != nil {
			log.Println(err, message)
			continue
		}
		err = c.database.DeleteOrder(order.ID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
