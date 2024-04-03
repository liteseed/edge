package contracts

import (
	"log"

	"github.com/everFinance/goar/types"
)

func (c *Context) UpdateStatus(dataItemId string, transactionId string) error {
	mId, err := c.ao.SendMessage(PROCESS, "", []types.Tag{{Name: "Action", Value: "Update"}, {Name: "DataItemId", Value: dataItemId}, {Name: "TransactionId", Value: transactionId}}, "", c.signer)
	if err != nil {
		return err
	}

	res, err := c.ao.ReadResult(PROCESS, mId)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
