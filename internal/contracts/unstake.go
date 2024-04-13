package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) Unstake() error {
	_, err := c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Unstake"}}, "", c.signer)
	return err
}
