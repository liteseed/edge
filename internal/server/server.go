package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/sdk-go/contract"
)

const (
	CONTENT_TYPE_OCTET_STREAM = "application/octet-stream"
	MAX_DATA_ITEM_SIZE        = uint(2 * 1024 * 1024 * 1024)
)

type Server struct {
	client   *client.Client
	contract *contract.Contract
	database *database.Config
	server   *http.Server
	store    *store.Store
	signer   *signer.Signer
}

func New(port string, version string, options ...func(*Server)) (*Server, error) {
	s := &Server{}
	for _, o := range options {
		o(s)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/", s.StatusGet(version))
	engine.POST("/tx", s.DataItemPost)
	engine.GET("/tx/:id", s.DataItemGet)
	engine.PUT("/tx/:id/:transaction_id", s.DataItemPut)
	engine.GET("/tx/:id/status", s.DataItemStatusGet)

	s.server = &http.Server{
		Addr:    port,
		Handler: engine,
	}
	return s, nil
}


func WithClient(c *client.Client) func(*Server) {
	return func(srv *Server) {
		srv.client = c
	}
}

func WithContracts(c *contract.Contract) func(*Server) {
	return func(srv *Server) {
		srv.contract = c
	}
}

func WithDatabase(db *database.Config) func(*Server) {
	return func(srv *Server) {
		srv.database = db
	}
}

func WithSigner(s *signer.Signer) func(*Server) {
	return func(srv *Server) {
		srv.signer = s
	}
}

func WithStore(s *store.Store) func(*Server) {
	return func(srv *Server) {
		srv.store = s
	}
}

func (srv *Server) Start() error {
	return srv.server.ListenAndServe()
}

func (srv *Server) Shutdown() error {
	return srv.server.Shutdown(context.TODO())
}
