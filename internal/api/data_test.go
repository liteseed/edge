package api

import (
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPostDataParseError(t *testing.T) {
	app, api := NewApiTest()
	app.POST("/data", api.PostData)
	r := PerformRequestWithBody(app, "POST", "/data", `{}`)
	assert.Assert(t, http.StatusBadRequest, r.Code)
}

func TestPostDataOk(t *testing.T) {
}
