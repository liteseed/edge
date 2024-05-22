package server

import (
	"context"
	"net/http"

	"github.com/everFinance/goar"
	"github.com/gin-gonic/gin"
	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
)

const (
	CONTENT_TYPE_OCTET_STREAM = "application/octet-stream"
	MAX_DATA_ITEM_SIZE        = 2 * 1024 * 1024 * 1024
)

type Server struct {
	contract   *contracts.Context
	database   *database.Config
	gatewayUrl string
	server     *http.Server
	store      *store.Store
	wallet     *goar.Wallet
}

func New(port string, version string, gatewayUrl string, options ...func(*Server)) (*Server, error) {
	s := &Server{}
	for _, o := range options {
		o(s)
	}
	s.gatewayUrl = gatewayUrl
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/", s.StatusGet(version))
	engine.POST("/tx", s.DataItemPost)

	s.server = &http.Server{
		Addr:    port,
		Handler: engine,
	}
	return s, nil
}

func WithContracts(contract *contracts.Context) func(*Server) {
	return func(c *Server) {
		c.contract = contract
	}
}

func WithDatabase(db *database.Config) func(*Server) {
	return func(c *Server) {
		c.database = db
	}
}

func WithStore(s *store.Store) func(*Server) {
	return func(c *Server) {
		c.store = s
	}
}
func WithWallet(w *goar.Wallet) func(*Server) {
	return func(c *Server) {
		c.wallet = w
	}
}
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}
