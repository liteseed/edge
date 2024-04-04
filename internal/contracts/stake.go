package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) Stake(url string) error {
	mId, err := c.ao.SendMessage(PROCESS, "", []types.Tag{{Name: "Action", Value: "Stake"}, {Name: "Url", Value: url}}, "", c.signer)
	if err != nil {
		return err
	}

	_, err = c.ao.ReadResult(PROCESS, mId)
	if err != nil {
		return err
	}

	return nil
}
