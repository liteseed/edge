package cron

import (
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/liteseed/edge/internal/database/schema"
	"github.com/liteseed/goar/wallet"
	"github.com/liteseed/sdk-go/contract"
	"github.com/stretchr/testify/assert"
)

const TEST_PROCESS = "PWSr59Cf6jxY7aA_cfz69rs0IiJWWbmQA8bAKknHeMo"

func mint(address string) {
	_, err := http.Get("http://localhost:1984/mint/" + address + "/1000000000000")
	if err != nil {
		panic(0)
	}
	mine()
}

func mine() {
	_, err := http.Get("http://localhost:1984/mine")
	if err != nil {
		panic(0)
	}
}

func TestCheckBundleConfirmations(t *testing.T) {
	data := []byte{1, 2, 3}

	user, err := wallet.New("http://localhost:1984")
	assert.NoError(t, err)
	mint(user.Signer.Address)
	mine()

	bundler, err := wallet.New("http://localhost:1984")
	assert.NoError(t, err)

	c := contract.New(TEST_PROCESS, bundler.Signer)

	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	crn, err := New(WithContracts(c), WithLogger(l), WithWallet(bundler))
	assert.NoError(t, err)

	p, err := user.Client.GetTransactionPrice(len(data), "")
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		tx := user.CreateTransaction(data, bundler.Signer.Address, p, nil)

		_, err = user.SignTransaction(tx)
		assert.NoError(t, err)

		err = user.SendTransaction(tx)
		assert.NoError(t, err)
		for i := 0; i < 11; i++ {
			mine()
		}

		u := crn.checkBundleConfirmations(tx.ID)
		assert.Equal(t, schema.Status("queued"), u.Status)
	})

	t.Run("Not Enough Confirmation", func(t *testing.T) {
		tx := user.CreateTransaction(data, bundler.Signer.Address, p, nil)

		_, err = user.SignTransaction(tx)
		assert.NoError(t, err)

		err = user.SendTransaction(tx)
		assert.NoError(t, err)

		for i := 0; i < 7; i++ {
			mine()
		}

		u := crn.checkBundleConfirmations(tx.ID)
		assert.Nil(t, u)
	})
}
