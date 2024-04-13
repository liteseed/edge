package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) GetBalance() (string, error) {
	mId, err := c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Balance"}}, "", c.signer)
	if err != nil {
		return "", err
	}

	result, err := c.ao.ReadResult(c.process, mId)
	if err != nil {
		return "", err
	}

	return result.Messages[0]["Data"].(string), nil
}
