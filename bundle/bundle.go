package bundle

import (
	"log"

	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
)

// Take raw data and generate a BundleItem
func PackBundleItem(data []byte) *types.BundleItem {
	bundleItem, err := utils.DecodeBundleItem(data)
	if err != nil {
		log.Fatal(err)
	}
	return bundleItem
}

