package cron

import (
	"github.com/everFinance/goar"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/robfig/cron/v3"
)

type Context struct {
	C        *cron.Cron
	database *database.Context
	store    *store.Store
	wallet   *goar.Wallet
}

func New(options ...func(*Context)) (*Context, error) {
	c := &Context{C: cron.New()}
	for _, o := range options {
		o(c)
	}
	return c, nil
}

func WithDatabase(db *database.Context) func(*Context) {
	return func(c *Context) {
		c.database = db
	}
}

func WithWallet(s *goar.Wallet) func(*Context) {
	return func(c *Context) {
		c.wallet = s
	}
}

func WithStore(s *store.Store) func(*Context) {
	return func(c *Context) {
		c.store = s
	}
}

func (c *Context) Start() {
	c.C.Start()
}

func (c *Context) Stop() {
	c.C.Stop()
}

func (c *Context) Add(spec string) error {
	_, err := c.C.AddFunc(spec, c.postBundle)
	return err
}
