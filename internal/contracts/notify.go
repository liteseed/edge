package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) Notify(dataItemId string, transactionId string) error {
	_, err := c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Notify"}, {Name: "DataItemId", Value: dataItemId}, {Name: "TransactionId", Value: transactionId}}, "", c.signer)
	return err
}
