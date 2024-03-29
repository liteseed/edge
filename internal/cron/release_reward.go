package cron

import (
	"log"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Context) ReleaseReward() {
	o, err := c.database.GetOrdersByStatus(schema.Reward)
	if err != nil {
		log.Println(err)
		return
	}

	for _, order := range *o {
		message, err := c.ao.SendMessage(PROCESS, "release", []transaction.Tag{{Name: "Action", Value: "Release"}, {Name: "Transaction", Value: order.ID}}, "", &signer.Signer{S: c.wallet.Signer})
		if err != nil {
			log.Println(err, message)
			continue
		}
	}
}
