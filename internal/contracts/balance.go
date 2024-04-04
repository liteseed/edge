package contracts

import (
	"github.com/everFinance/goar/types"
)

func (c *Context) GetBalance() (string, error) {
	mId, err := c.ao.SendMessage(PROCESS, "", []types.Tag{{Name: "Action", Value: "Balance"}}, "", c.signer)
	if err != nil {
		return "", err
	}

	result, err := c.ao.ReadResult(PROCESS, mId)
	if err != nil {
		return "", err
	}

	return result.Messages[0]["Data"].(string), nil
}
