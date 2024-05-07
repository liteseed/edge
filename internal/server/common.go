package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"strconv"

	"github.com/everFinance/goar/types"
	"github.com/liteseed/edge/internal/contracts"
)

func verify(contract *contracts.Context, dataItem *types.BundleItem) (bool, error) {
	res, err := contract.GetUpload(dataItem.Id)
	if err != nil {
		return false, err
	}
	if res == nil {
		return false, errors.New("no upload")
	}
	rawData, err := base64.RawURLEncoding.DecodeString(dataItem.Data)
	if err != nil {
		return false, err
	}

	price, err := getPrice(len(rawData))
	if err != nil {
		return false, err
	}

	log.Println(price)
	// if quantity >= 0 {
	// 	return false, errors.New("quantity isn't enough")
	// }

	return true, nil
}

func getPrice(b int) (int, error) {

	res, err := http.Get(fmt.Sprint("https://arweave.net/price/", b))
	if err != nil {
		return 0, err
	}

	r, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil
	}

	i, err := strconv.Atoi(string(r))
	if err != nil {
		return 0, nil
	}

	return i, nil
}
