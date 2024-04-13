package contracts

import (
	"log"

	"github.com/everFinance/goar/types"
)

func (c *Context) Release(dataItemId string) error {
	message, err := c.ao.SendMessage(c.process, "release", []types.Tag{{Name: "Action", Value: "Release"}, {Name: "DataItemId", Value: dataItemId}}, "", c.signer)
	if err != nil {
		return err
	}
	log.Println(message)
	return nil
}
