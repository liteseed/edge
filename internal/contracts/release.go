package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) Release(dataItemId string) error {
	_, err := c.ao.SendMessage(c.process, "release", []types.Tag{{Name: "Action", Value: "Release"}, {Name: "DataItemId", Value: dataItemId}}, "", c.signer)
	return err
}
