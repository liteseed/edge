package cron

import (
	"log"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/liteseed/edge/internal/database/schema"
)

func (c *Context) ReleaseReward() {
	o, err := c.database.GetOrdersByStatus(schema.Reward)
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
		message, err := c.ao.SendMessage(PROCESS, "release", []types.Tag{{Name: "Action", Value: "Release"}, {Name: "Transaction", Value: order.ID}}, "", itemSigner)
		if err != nil {
			log.Println(err, message)
			continue
		}
	}
}
