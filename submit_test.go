package bungo

import (
	"os"
	"testing"

	"github.com/everFinance/goar"
	"github.com/liteseed/bungo/store"
	"github.com/stretchr/testify/assert"
)

func TestFetchAndStoreTx(t *testing.T) {
	arId := "O3VwBusl0PNKusWcDF44uPt-sNuhywgeKxOmQpDqGc0"
	dbPath := "./data/tmp.db"
	arNode := "https://arweave.net"
	cli := goar.NewClient(arNode)
	s, err := store.NewBoltStore(dbPath)
	assert.NoError(t, err)
	aa := &Bungo{store: s, arCli: cli}
	err = aa.FetchAndStoreTx(arId)
	assert.NoError(t, err)
	err = os.RemoveAll(dbPath)
	assert.NoError(t, err)
}

func TestSaveSubmitTx(t *testing.T) {
	arId := "O3VwBusl0PNKusWcDF44uPt-sNuhywgeKxOmQpDqGc0"
	dbPath := "./data/tmp.db"
	arNode := "https://arweave.net"
	cli := goar.NewClient(arNode)
	s, err := store.NewBoltStore(dbPath)
	assert.NoError(t, err)
	aa := &Bungo{store: s, arCli: cli}
	tx, err := cli.GetTransactionByID(arId)
	assert.NoError(t, err)
	err = aa.SaveSubmitTx(*tx)
	assert.NoError(t, err)
	err = os.RemoveAll(dbPath)
	assert.NoError(t, err)
}

func TestSyncAddTxDataEndOffset(t *testing.T) {
	arId := "O3VwBusl0PNKusWcDF44uPt-sNuhywgeKxOmQpDqGc0"
	dbPath := "./data/tmp.db"
	arNode := "https://arweave.net"
	cli := goar.NewClient(arNode)
	s, err := store.NewBoltStore(dbPath)
	assert.NoError(t, err)
	aa := &Bungo{store: s, arCli: cli}
	tx, err := cli.GetTransactionByID(arId)
	assert.NoError(t, err)
	err = aa.syncAddTxDataEndOffset(tx.DataRoot, tx.DataSize)
	assert.NoError(t, err)
	err = os.RemoveAll(dbPath)
	assert.NoError(t, err)
}
