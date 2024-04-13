package contracts

import (
	"log"

	"github.com/everFinance/goar/types"
)

func (c *Context) Notify(dataItemId string, transactionId string) error {
	mId, err := c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Notify"}, {Name: "DataItemId", Value: dataItemId}, {Name: "TransactionId", Value: transactionId}}, "", c.signer)
	if err != nil {
		return err
	}

	res, err := c.ao.ReadResult(c.process, mId)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
