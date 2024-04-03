package contracts

import (
	"encoding/json"

	"github.com/everFinance/goar/types"
)

type GetUploadResponse struct {
	Status   string `json:"status"`
	Quantity string `json:"quantity"`
	Block    string `json:"block"`
	Bundler  string `json:"bundler"`
}

func (c *Context) GetUpload(dataItemId string) (*GetUploadResponse, error) {
	mId, err := c.ao.SendMessage(PROCESS, "", []types.Tag{{Name: "Action", Value: "Upload"}, {Name: "DataItemId", Value: dataItemId}, {Name: "Status", Value: ""}}, "", c.signer)
	if err != nil {
		return nil, err
	}

	result, err := c.ao.ReadResult(PROCESS, mId)
	if err != nil {
		return nil, err
	}

	var response GetUploadResponse
	err = json.Unmarshal([]byte(result.Messages[0]["Data"].(string)), &response)
	if err != nil {
		return nil, err
	}
	
	return &response, err
}
