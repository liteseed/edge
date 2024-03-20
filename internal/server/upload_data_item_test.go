package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/database/schema"
	"github.com/liteseed/edge/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestUploadDataItem(t *testing.T) {
	defer os.RemoveAll("./temp-upload-data-item")
	t.Parallel()
	database, _ := database.New("sqlite", "./temp-upload-data-item/sqlite")
	signer, _ := signer.New("../../data/signer.json")
	store := store.New("pebble", "./temp-upload-data-item/pebble")
	data, _ := os.ReadFile("../../test/1115BDataItem")

	gin.SetMode(gin.TestMode)
	server := New(database, signer, store)

	t.Run("content-type:application/json, content-length:1115, data:1115BDataItem", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/json", "content-length": "1115"}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data-item", bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "content-type: unsupported"}, resErr)
	})

	t.Run("content-type:application/octet-stream, content-length:0,  data:1115BDataItem", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "0"}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data-item", bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "content-length: supported range 1B - 1073824B"}, resErr)
	})
	t.Run("content-type:application/octet-stream, content-length:1115, data:nil", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "1115"}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data-item", nil)
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "body: required"}, resErr)
	})
	t.Run("content-type:application/octet-stream, content-length:4, data:1115BDataItem", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "4"}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data-item", bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resErr responseError
		err := json.Unmarshal(w.Body.Bytes(), &resErr)

		assert.NoError(t, err)
		assert.Equal(t, responseError{Error: "content-length, body: length mismatch (4, 1115)"}, resErr)
	})
	t.Run("content-type:application/octet-stream, content-length:1115, data:data-item", func(t *testing.T) {
		headers := map[string]string{"content-type": "application/octet-stream", "content-length": "1115"}
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/data-item", bytes.NewReader(data))
		req.Header.Set("content-type", headers["content-type"])
		req.Header.Set("content-length", headers["content-length"])
		server.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var res UploadDataResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Id)

		o, err := database.GetOrder(uuid.MustParse(res.Id))
		assert.NoError(t, err)

		status, err := o.Status.Value()
		assert.NoError(t, err)
		assert.Equal(t, schema.Queued, status)

		rawData, err := store.Get(o.StoreID)
		assert.NoError(t, err)

		dataItem, err := transaction.DecodeDataItem(rawData)
		assert.NoError(t, err)
		assert.ElementsMatch(t, data, dataItem.Raw)
	})
}
