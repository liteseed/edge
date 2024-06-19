package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/liteseed/aogo"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/wallet"
	"github.com/liteseed/sdk-go/contract"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
)

func TestNewServer(t *testing.T) {
	srv, err := New(":8080", "test")
	assert.NoError(t, err)
	assert.NotNil(t, srv)
}

func TestStatusHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		arweave := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
			_, err := writer.Write([]byte(`{"network":"arweave.N.1","version": 5,"release": 69,"height": 1447908,"current":"XwcWkjKLbXlDg8QcagmW0AN6c2V3y0lyHEaPLT2tUf8vH9kKM5OlfYmfKQtd6XxI","blocks": 1447909,"peers": 307,"queue_length": 0,"node_stat)e_latency": 1}`))
			assert.NoError(t, err)
		}))
		defer arweave.Close()
		w, err := wallet.FromPath("../../test/signer.json", arweave.URL)
		assert.NoError(t, err)
		srv, _ := New(":8080", "test", WithWallet(w))

		rcd := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		srv.server.Handler.ServeHTTP(rcd, req)

		assert.Equal(t, http.StatusOK, rcd.Code)
		assert.Equal(t, fmt.Sprintf(`{"Address":"3XTR7MsJUD9LoaiFRdWswzX1X5BR7AQdl1x2v2zIVck","Gateway":{"Block-Height":1447908,"Status":"ok","URL":"%s"},"Name":"Edge","Version":"test"}`, arweave.URL), rcd.Body.String())

	})

	t.Run("Gateway Error", func(t *testing.T) {
		arweave := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) { writer.WriteHeader(http.StatusNotFound) }))
		defer arweave.Close()
		w, err := wallet.FromPath("../../test/signer.json", arweave.URL)
		assert.NoError(t, err)
		srv, _ := New(":8080", "test", WithWallet(w))

		rcd := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		srv.server.Handler.ServeHTTP(rcd, req)

		assert.Equal(t, http.StatusFailedDependency, rcd.Code)
		assert.Equal(t, fmt.Sprintf(`{"Address":"3XTR7MsJUD9LoaiFRdWswzX1X5BR7AQdl1x2v2zIVck","Gateway":{"Block-Height":-1,"Status":"failed","URL":"%s"},"Name":"Edge","Version":"test"}`, arweave.URL), rcd.Body.String())

	})
}

func TestDataItemPost(t *testing.T) {
	arweave := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		_, err := writer.Write([]byte(`{"network":"arweave.N.1","version": 5,"release": 69,"height": 1447908,"current":"XwcWkjKLbXlDg8QcagmW0AN6c2V3y0lyHEaPLT2tUf8vH9kKM5OlfYmfKQtd6XxI","blocks": 1447909,"peers": 307,"queue_length": 0,"node_stat)e_latency": 1}`))
		assert.NoError(t, err)
	}))

	w, err := wallet.FromPath("../../test/signer.json", arweave.URL)
	assert.NoError(t, err)

	d := w.CreateDataItem([]byte{1, 2, 3}, "", "", []tag.Tag{})
	_, err = w.SignDataItem(d)
	assert.NoError(t, err)

	s := store.New("../../test/store")

	mockDb, mock, _ := sqlmock.New()
	db, err := database.FromDialector(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	}))

	assert.NoError(t, err)

	srv, err := New(":8000", "test", WithDatabase(db), WithStore(s), WithWallet(w))
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "orders" ("id","status","transaction_id","bundle_id","size","deadline_height") VALUES ($1,$2,$3,$4,$5,$6)`)).WithArgs(d.ID, "created", "", "", 1047, 1448108).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		rcd := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer(d.Raw))
		req.Header.Set("content-type", "application/octet-stream")
		req.Header.Set("content-length", strconv.Itoa(len(d.Raw)))

		srv.server.Handler.ServeHTTP(rcd, req)

		assert.Equal(t, http.StatusCreated, rcd.Code)
		assert.Equal(t, fmt.Sprintf("{\"id\":\"%s\",\"owner\":\"3XTR7MsJUD9LoaiFRdWswzX1X5BR7AQdl1x2v2zIVck\",\"dataCaches\":[\"%s\"],\"deadlineHeight\":1448108,\"fastFinalityIndexes\":[\"%s\"],\"version\":\"1.0.0\"}", d.ID, arweave.URL, arweave.URL), rcd.Body.String())
	})

	t.Run("Missing", func(t *testing.T) {
		rcd := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer(d.Raw))
		srv.server.Handler.ServeHTTP(rcd, req)
		assert.Equal(t, http.StatusBadRequest, rcd.Code)
		assert.Equal(t, `{"error":"required header(s) - content-type, content-length"}`, rcd.Body.String())
	})
	t.Run("Invalid", func(t *testing.T) {
		rcd := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer(d.Raw))
		req.Header.Set("content-type", "application/json")
		req.Header.Set("content-length", strconv.Itoa(len(d.Raw)))
		srv.server.Handler.ServeHTTP(rcd, req)

		assert.Equal(t, http.StatusBadRequest, rcd.Code)
		assert.Equal(t, `{"error":"required header(s) - content-type: application/octet-stream"}`, rcd.Body.String())
	})
	t.Run("Invalid Content Type", func(t *testing.T) {
		rcd := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer(d.Raw))
		req.Header.Set("content-type", "application/json")
		req.Header.Set("content-length", strconv.Itoa(len(d.Raw)))
		srv.server.Handler.ServeHTTP(rcd, req)

		assert.Equal(t, http.StatusBadRequest, rcd.Code)
		assert.Equal(t, `{"error":"required header(s) - content-type: application/octet-stream"}`, rcd.Body.String())
	})
	t.Run("Invalid Content Length", func(t *testing.T) {
		rcd := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer(d.Raw))
		req.Header.Set("content-type", "application/octet-stream")
		req.Header.Set("content-length", "-100")
		srv.server.Handler.ServeHTTP(rcd, req)

		assert.Equal(t, http.StatusBadRequest, rcd.Code)
		assert.Equal(t, `{"error":"content-length, body: length mismatch (-100, 1047)"}`, rcd.Body.String())
	})

	t.Run("Nil Body", func(t *testing.T) {
		rcd := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tx", nil)
		req.Header.Set("content-type", "application/octet-stream")
		req.Header.Set("content-length", strconv.Itoa(len(d.Raw)))
		srv.server.Handler.ServeHTTP(rcd, req)
		assert.Equal(t, http.StatusBadRequest, rcd.Code)
		assert.Equal(t, `{"error":"cannot read nil body"}`, rcd.Body.String())
	})

	t.Run("Invalid Body", func(t *testing.T) {
		rcd := httptest.NewRecorder()

		req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer([]byte{1, 2, 3}))
		req.Header.Set("content-type", "application/octet-stream")
		req.Header.Set("content-length", strconv.Itoa(len(d.Raw)))
		srv.server.Handler.ServeHTTP(rcd, req)
		assert.Equal(t, http.StatusBadRequest, rcd.Code)
		assert.Equal(t, `{"error":"content-length, body: length mismatch (1047, 3)"}`, rcd.Body.String())
	})

	t.Run("Invalid Data Item", func(t *testing.T) {
		rcd := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer([]byte{1, 2, 3}))
		req.Header.Set("content-type", "application/octet-stream")
		req.Header.Set("content-length", "3")
		srv.server.Handler.ServeHTTP(rcd, req)
		assert.Equal(t, http.StatusBadRequest, rcd.Code)
		assert.Equal(t, `{"error":"failed to decode data item"}`, rcd.Body.String())
	})

}

func TestDataItemPut(t *testing.T) {
	w, err := wallet.FromPath("../../test/signer.json", "")
	assert.NoError(t, err)

	mu := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"id":"id", "message": ""}`))
		assert.NoError(t, err)
	}))
	defer mu.Close()

	ao, err := aogo.New(aogo.WthMU(mu.URL))
	assert.NoError(t, err)

	c := contract.Custom(ao, "process", w.Signer)
	mockDb, mock, _ := sqlmock.New()
	db, err := database.FromDialector(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	}))

	assert.NoError(t, err)

	srv, err := New(":8000", "test", WithDatabase(db), WithContracts(c))
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "orders" SET "status"=$1,"transaction_id"=$2 WHERE id = $3`)).WithArgs("queued", "payment", "dataitem").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		rcd := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tx/dataitem/payment", nil)

		srv.server.Handler.ServeHTTP(rcd, req)

		assert.Equal(t, http.StatusAccepted, rcd.Code)
		assert.Equal(t, "{\"id\":\"dataitem\",\"payment_id\":\"payment\"}", rcd.Body.String())
	})
}
