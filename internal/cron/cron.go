package cron

import (
	"github.com/liteseed/argo/signer"
	"github.com/liteseed/bungo/internal/database"
	"github.com/liteseed/bungo/internal/store"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	C        *cron.Cron
	database *database.Database
	store    *store.Store
	signer   *signer.Signer
}

func New(options ...func(*Cron)) (*Cron, error) {
	c := &Cron{C: cron.New()}
	for _, o := range options {
		o(c)
	}
	return c, nil
}

func WithDatabase(db *database.Database) func(*Cron) {
	return func(c *Cron) {
		c.database = db
	}
}

func WithSigner(s *signer.Signer) func(*Cron) {
	return func(c *Cron) {
		c.signer = s
	}
}

func WithStore(s *store.Store) func(*Cron) {
	return func(c *Cron) {
		c.store = s
	}
}

func (c *Cron) Start() {
	c.C.Start()
}

func (c *Cron) Stop() {
	c.C.Stop()
}

func (c *Cron) Add(spec string) error {
	_, err := c.C.AddFunc(spec, c.postBundle)
	return err
}
