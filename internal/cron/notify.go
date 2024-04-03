package cron

import (
	"log"

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
		err := c.contract.UpdateStatus(order.ID, order.TransactionID)
		if err != nil {
			log.Println(err)
			continue
		}
		err = c.database.UpdateStatus(order.ID, schema.Reward)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
