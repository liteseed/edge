package cron

import (
	"log/slog"

	"github.com/everFinance/goar"
	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/robfig/cron/v3"
)

type Context struct {
	C        *cron.Cron
	contract *contracts.Context
	database *database.Context
	logger   *slog.Logger
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

func WthContracts(contract *contracts.Context) func(*Context) {
	return func(c *Context) {
		c.contract = contract
	}
}

func WithDatabase(db *database.Context) func(*Context) {
	return func(c *Context) {
		c.database = db
	}
}

func WithLogger(logger *slog.Logger) func(*Context) {
	return func(c *Context) {
		c.logger = logger
	}
}

func WithStore(s *store.Store) func(*Context) {
	return func(c *Context) {
		c.store = s
	}
}
func WithWallet(s *goar.Wallet) func(*Context) {
	return func(c *Context) {
		c.wallet = s
	}
}

func (c *Context) Start() {
	c.C.Start()
}

func (c *Context) Stop() {
	c.C.Stop()
}

func (c *Context) PostBundle(spec string) error {
	_, err := c.C.AddFunc(spec, c.postBundle)
	if err != nil {
		return err
	}
	_, err = c.C.AddFunc(spec, c.notify)
	if err != nil {
		return err
	}
	_, err = c.C.AddFunc(spec, c.SyncStatus)
	if err != nil {
		return err
	}
	_, err = c.C.AddFunc(spec, c.ReleaseReward)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) Notify() error {
	_, err := c.C.AddFunc("* * * * *", c.notify)
	return err
}
