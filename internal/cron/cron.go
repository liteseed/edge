package cron

import (
	"log/slog"

	"github.com/everFinance/goar"
	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	c        *cron.Cron
	contract *contracts.Context
	database *database.Config
	logger   *slog.Logger
	store    *store.Store
	wallet   *goar.Wallet
}

type Option = func(*Cron)

func New(options ...func(*Cron)) (*Cron, error) {
	c := &Cron{c: cron.New()}
	for _, o := range options {
		o(c)
	}
	return c, nil
}

func WthContracts(contract *contracts.Context) Option {
	return func(c *Cron) {
		c.contract = contract
	}
}

func WithDatabase(db *database.Config) Option {
	return func(c *Cron) {
		c.database = db
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Cron) {
		c.logger = logger
	}
}

func WithStore(s *store.Store) Option {
	return func(c *Cron) {
		c.store = s
	}
}
func WithWallet(s *goar.Wallet) Option {
	return func(c *Cron) {
		c.wallet = s
	}
}

func (c *Cron) Start() {
	c.c.Start()
}

func (c *Cron) Shutdown() {
	c.c.Stop()
}

func (c *Cron) Setup(spec string) error {
	_, err := c.c.AddFunc(spec, c.PostBundle)
	if err != nil {
		return err
	}
	_, err = c.c.AddFunc(spec, c.Notify)
	if err != nil {
		return err
	}
	_, err = c.c.AddFunc(spec, c.SyncStatus)
	if err != nil {
		return err
	}
	_, err = c.c.AddFunc(spec, c.ReleaseReward)
	if err != nil {
		return err
	}
	return nil
}
