package ao

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/liteseed/bungo/internal/types"
)

// TODO: Replace with complete golang ao-connect implementation

var URL = os.Getenv("AO_CONNECT_URL")
var PROCESS = os.Getenv("NETWORK_CONTRACT_ID")


type SendMessageArgs struct {
	Data string `json:"data"`
	Tags []types.Tag  `json:"tags"`
}


func SendMessage(args SendMessageArgs) (string, error) {
	reqBody, err := json.Marshal(args)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(URL+"/"+PROCESS, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	messageId, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(messageId), nil
}

type ReadResultArgs struct {
	Message string `json:"message"`
}

type ReadResultResponse struct {
	Messages []any
	Spawns   []any
	Outputs  []any
	Errors   any
	GasUsed  int
}

func ReadResult(args ReadResultArgs) (*ReadResultResponse, error) {
	resp, err := http.Get(URL + "/" + PROCESS + "/" + args.Message)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &ReadResultResponse{}
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
