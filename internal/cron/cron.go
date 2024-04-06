package cron

import (
	"log/slog"

	"github.com/everFinance/goar"
	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/robfig/cron/v3"
)

type Config struct {
	c        *cron.Cron
	contract *contracts.Context
	database *database.Config
	logger   *slog.Logger
	store    *store.Store
	wallet   *goar.Wallet
}

type Option = func(*Config)

func New(options ...func(*Config)) (*Config, error) {
	c := &Config{c: cron.New()}
	for _, o := range options {
		o(c)
	}
	return c, nil
}

func WthContracts(contract *contracts.Context) Option {
	return func(c *Config) {
		c.contract = contract
	}
}

func WithDatabase(db *database.Config) Option {
	return func(c *Config) {
		c.database = db
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Config) {
		c.logger = logger
	}
}

func WithStore(s *store.Store) Option {
	return func(c *Config) {
		c.store = s
	}
}
func WithWallet(s *goar.Wallet) Option {
	return func(c *Config) {
		c.wallet = s
	}
}

func (c *Config) Start() {
	c.c.Start()
}

func (c *Config) Shutdown() {
	c.c.Stop()
}

func (c *Config) Setup(spec string) error {
	_, err := c.c.AddFunc(spec, c.postBundle)
	if err != nil {
		return err
	}
	_, err = c.c.AddFunc(spec, c.notify)
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
