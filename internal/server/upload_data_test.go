package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/everFinance/goar"
	"github.com/gin-gonic/gin"
	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
	"github.com/liteseed/edge/internal/store"
	"github.com/stretchr/testify/assert"
)

type responseError struct {
	Error string `json:"error"`
}

func TestUploadData(t *testing.T) {
	defer os.RemoveAll("./temp-upload-data")

	_ = os.Mkdir("./temp-upload-data", os.ModePerm)
	id := "AVASWERFDHTRE"
	database, _ := database.New("sqlite", "./temp-upload-data/sqlite")
	store := store.New("pebble", "./temp-upload-data/pebble")
	signer, _ := goar.NewSignerFromPath("../../data/signer.json")
	gin.SetMode(gin.TestMode)
	server := New(database, signer, store)

	t.Parallel()

	t.Run("content-type:application/json, content-length:3, data:{1,2,3}", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/json", "content-length": "3"}
		data := []byte{1, 2, 3}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data/"+id, bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "content-type: unsupported"}, resErr)
	})

	t.Run("content-type:application/octet-stream, content-length:0, data:{1,2,3}", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "0"}
		data := []byte{1, 2, 3}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data/"+id, bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "content-length: supported range 1B - 1073824B"}, resErr)
	})
	t.Run("content-type:application/octet-stream, content-length:3, data:nil", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "3"}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data/"+id, nil)
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "body: required"}, resErr)
	})
	t.Run("content-type:application/octet-stream, content-length:4, data:{1,2,3}", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "4"}
		data := []byte{1, 2, 3}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data/"+id, bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "content-length, body: length mismatch (4, 3)"}, resErr)
	})
	t.Run("content-type:application/octet-stream, content-length:3, data:{1,2,3}", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "3"}
		data := []byte{1, 2, 3}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data/"+id, bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		o, err := database.GetOrder(id)
		assert.NoError(t, err)

		status, err := o.Status.Value()
		assert.NoError(t, err)
		assert.Equal(t, schema.Queued, status)

		rawData, err := store.Get(o.StoreID)
		assert.NoError(t, err)

		dataItem, err := transaction.DecodeDataItem(rawData)
		assert.NoError(t, err)
		assert.Equal(t, base64.RawURLEncoding.EncodeToString(data), dataItem.Data)
	})
}
