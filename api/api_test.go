package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/bungo/cache"
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/store"
)

// NewApiTest returns new API test helper.
func NewApiTest() (*gin.Engine, *API) {
	gin.SetMode(gin.TestMode)

	cache, err := cache.NewBigCache(60 * time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewSqliteDatabase("./tmp/sqlite")
	if err = db.Migrate(); err != nil {
		log.Fatal(err)
	}

	store, err := store.NewBoltStore("./tmp/sqlite")
	if err != nil {
		log.Fatal(err)
	}
	a := New(cache, db, store)

	r := gin.Default()

	r.GET("/status", a.GetStatus)
	r.POST("/data", a.PostData)

	return r, a
}

// PerformRequest runs an API request with an empty request body.
// See https://medium.com/@craigchilds94/testing-gin-json-responses-1f258ce3b0b1
func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

// PerformRequestWithBody runs an API request with the request body as a string.
func PerformRequestWithBody(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	reader := strings.NewReader(body)
	req, _ := http.NewRequest(method, path, reader)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return w
}
