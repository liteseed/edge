package cron

import (
	"log"

	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Context) ReleaseReward() {
	o, err := c.database.GetOrdersByStatus(schema.Reward)
	if err != nil {
		log.Println(err)
		return
	}

	for _, order := range *o {
		err = c.contract.Release(order.ID)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
