package bundle

import (
	"os"
	"testing"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/liteseed/bungo/schema"
	"github.com/liteseed/bungo/store"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func setupBundle() (b *Bundle) {
	s, _ := store.NewBoltStore("./tmp.db")
	perFeeMap := map[string]schema.Fee{
		"AR": {
			Currency: "AR",
			Decimals: 0,
			Base:     decimal.New(5, 0),
			PerChunk: decimal.New(1, 0),
		},
	}
	b = &Bundle{
		bundlePerFeeMap: perFeeMap,
		store:           s,
	}
	return
}

func teardown() {
	os.RemoveAll("./tmp.db")
}

func TestBundleCalcItemFee(t *testing.T) {
	b := setupBundle()
	defer teardown()

	size0 := int64(types.MAX_CHUNK_SIZE - 1)
	size1 := int64(types.MAX_CHUNK_SIZE + 1)
	size2 := int64(types.MAX_CHUNK_SIZE)
	size3 := int64(0)

	res, err := b.CalcItemFee("AR", size0)
	assert.NoError(t, err)
	assert.Equal(t, "6", res.FinalFee)

	res, err = b.CalcItemFee("AR", size1)
	assert.NoError(t, err)
	assert.Equal(t, "7", res.FinalFee)

	res, err = b.CalcItemFee("AR", size2)
	assert.NoError(t, err)
	assert.Equal(t, "6", res.FinalFee)

	res, err = b.CalcItemFee("AR", size3)
	assert.NoError(t, err)
	assert.Equal(t, "5", res.FinalFee)

}

func TestSaveDelItem(t *testing.T) {
	b := setupBundle()
	defer teardown()

	signer, err := goar.NewSignerFromPath("../data/bundler-keyfile.json") // your key file path
	assert.NoError(t, err)
	bundlSigner, err := goar.NewItemSigner(signer)
	assert.NoError(t, err)
	item, err := bundlSigner.CreateAndSignItem([]byte("data 01"), "", "", nil)
	assert.NoError(t, err)

	err = b.store.AtomicSaveItem(item)
	assert.NoError(t, err)
	err = b.store.AtomicDeleteItem(item.Id)
	assert.NoError(t, err)
}

func TestParseAndSaveBundleItems(t *testing.T) {
	b := setupBundle()
	defer teardown()
	arId := "p5lopWlVbGvBiPy5kevaEWtwHHpl0coUeUv9qNQh6NA"
	node := "https://arweave.net"

	client := goar.NewClient(node)
	data, err := client.GetTransactionData(arId)
	assert.NoError(t, err)
	err = b.ParseAndSaveBundleItems(arId, data)
	assert.NoError(t, err)
}
