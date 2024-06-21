package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/liteseed/goar/transaction/data_item"
	"github.com/liteseed/goar/wallet"
	"gorm.io/driver/postgres"
)

func DataItem() *data_item.DataItem {
	b, _ := os.ReadFile("../../test/1115BDataItem")
	dataItem, _ := data_item.Decode(b)
	return dataItem
}

func Gateway() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := strings.Split(r.URL.Path, "/")
			if r.URL.Path == "/" || r.URL.Path == "/info" {
				_, _ = w.Write([]byte(`{"network":"arweave.N.1","version": 5,"release": 69,"height": 1447908,"current":"XwcWkjKLbXlDg8QcagmW0AN6c2V3y0lyHEaPLT2tUf8vH9kKM5OlfYmfKQtd6XxI","blocks": 1447909,"peers": 307,"queue_length": 0,"node_stat)e_latency": 1}`))
			} else if p[1] == "price" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("10000"))
			} else if p[1] == "tx_anchor" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(""))
			} else if p[len(p)-1] == "status" {
				w.WriteHeader(http.StatusOK)
				if p[len(p)-2] == "failbundle" {
					_, _ = w.Write([]byte(`{"block_height":1000,"block_indep_hash":"block_indep_hash","number_of_confirmations":11}`))
				} else {
					_, _ = w.Write([]byte(`{"block_height":1000,"block_indep_hash":"block_indep_hash","number_of_confirmations":26}`))
				}
			} else if p[1] == "tx" {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
}

func Database() (sqlmock.Sqlmock, *database.Database) {
	mockDb, mock, _ := sqlmock.New()
	db, _ := database.FromDialector(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	}))
	return mock, db
}

func Store() *store.Store {
	return store.New("../../test/store")
}

func Wallet(gateway string) *wallet.Wallet {
	w, _ := wallet.FromPath("../../test/signer.json", gateway)
	return w
}

func CU() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
}

func MU() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"id", "message": ""}`))
	}))
}
