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
	MAX_DATA_ITEM_SIZE        = 1_073_824
)

type Server struct {
	contract *contracts.Context
	database *database.Config
	server   *http.Server
	signer   *goar.Signer
	store    *store.Store
	logger   *slog.Logger
}

func New(port string, version string, options ...func(*Server)) (*Server, error) {

	s := &Server{}
	for _, o := range options {
		o(s)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(JSONLogMiddleware(s.logger))

	engine.GET("/", s.getStatus(version))
	engine.POST("/tx", s.uploadDataItem)

	s.server = &http.Server{
		Addr:    port,
		Handler: engine,
	}
	return s, nil
}

func WthContracts(contract *contracts.Context) func(*Server) {
	return func(c *Server) {
		c.contract = contract
	}
}

func WithDatabase(db *database.Config) func(*Server) {
	return func(c *Server) {
		c.database = db
	}
}

func WithLogger(logger *slog.Logger) func(*Server) {
	return func(c *Server) {
		c.logger = logger
	}
}

func WithStore(s *store.Store) func(*Server) {
	return func(c *Server) {
		c.store = s
	}
}
func WithWallet(s *goar.Signer) func(*Server) {
	return func(c *Server) {
		c.signer = s
	}
}
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}
