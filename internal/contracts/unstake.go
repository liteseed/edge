package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) Unstake() error {
	mId, err := c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Unstake"}}, "", c.signer)
	if err != nil {
		return err
	}

	_, err = c.ao.ReadResult(c.process, mId)
	if err != nil {
		return err
	}

	return nil
}
