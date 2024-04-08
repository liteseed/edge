package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/everFinance/goar"
	"github.com/gin-gonic/gin"
	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
)

const (
	CONTENT_TYPE_OCTET_STREAM = "application/octet-stream"
	MAX_DATA_SIZE             = 1_073_824
	MAX_DATA_ITEM_SIZE        = 1_073_824
)

type Config struct {
	contract *contracts.Context
	database *database.Config
	server   *http.Server
	signer   *goar.Signer
	store    *store.Store
	logger   *slog.Logger
}

func New(port string, options ...func(*Config)) (*Config, error) {

	s := &Config{}
	for _, o := range options {
		o(s)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(JSONLogMiddleware(s.logger))
	engine.GET("/status", s.Status)
	engine.POST("/tx", s.uploadDataItem)

	s.server = &http.Server{
		Addr:    port,
		Handler: engine,
	}
	return s, nil
}

func WthContracts(contract *contracts.Context) func(*Config) {
	return func(c *Config) {
		c.contract = contract
	}
}

func WithDatabase(db *database.Config) func(*Config) {
	return func(c *Config) {
		c.database = db
	}
}

func WithLogger(logger *slog.Logger) func(*Config) {
	return func(c *Config) {
		c.logger = logger
	}
}

func WithStore(s *store.Store) func(*Config) {
	return func(c *Config) {
		c.store = s
	}
}
func WithWallet(s *goar.Signer) func(*Config) {
	return func(c *Config) {
		c.signer = s
	}
}
func (s *Config) Start() error {
	return s.server.ListenAndServe()
}
func (s *Config) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}
