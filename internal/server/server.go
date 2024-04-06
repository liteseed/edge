package server

import (
	"log/slog"

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
	database *database.Context
	engine   *gin.Engine
	signer   *goar.Signer
	store    *store.Store
	logger   *slog.Logger
}

func New(options ...func(*Config)) (*Config, error) {
	s := &Config{engine: gin.New()}
	for _, o := range options {
		o(s)
	}
	return s, nil
}

func WthContracts(contract *contracts.Context) func(*Config) {
	return func(c *Config) {
		c.contract = contract
	}
}

func WithDatabase(db *database.Context) func(*Config) {
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

func (s *Config) Run(port string) error {
	s.engine = gin.New()

	s.engine.Use(gin.Recovery())
	s.engine.Use(JSONLogMiddleware(s.logger))
	s.engine.GET("/status", s.Status)
	return s.engine.Run(port)
}
