package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) Stake(url string) error {
	_, err := c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Stake"}, {Name: "Url", Value: url}}, "", c.signer)
	return err
}
