package cron

import (
	"log/slog"

	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/sdk-go/contract"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	c        *cron.Cron
	client   *client.Client
	contract *contract.Contract
	database *database.Config
	logger   *slog.Logger
	signer   *signer.Signer
	store    *store.Store
}

type Option = func(*Cron)

func New(options ...func(*Cron)) (*Cron, error) {
	c := &Cron{c: cron.New()}
	for _, o := range options {
		o(c)
	}
	return c, nil
}

func WithClient(c *client.Client) Option {
	return func(crn *Cron) {
		crn.client = c
	}
}

func WithContracts(contract *contract.Contract) Option {
	return func(crn *Cron) {
		crn.contract = contract
	}
}

func WithDatabase(db *database.Config) Option {
	return func(crn *Cron) {
		crn.database = db
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(crn *Cron) {
		crn.logger = logger
	}
}

func WithSigner(s *signer.Signer) Option {
	return func(crn *Cron) {
		crn.signer = s
	}
}

func WithStore(s *store.Store) Option {
	return func(crn *Cron) {
		crn.store = s
	}
}

func (crn *Cron) Start() {
	crn.c.Start()
}

func (crn *Cron) Shutdown() {
	crn.c.Stop()
}

func (crn *Cron) Setup(spec string) error {
	_, err := crn.c.AddFunc(spec, crn.CheckBundleConfirmation)
	if err != nil {
		return err
	}
	_, err = crn.c.AddFunc(spec, crn.PostBundle)
	if err != nil {
		return err
	}
	_, err = crn.c.AddFunc(spec, crn.Posted)
	if err != nil {
		return err
	}
	_, err = crn.c.AddFunc(spec, crn.Release)
	if err != nil {
		return err
	}
	_, err = crn.c.AddFunc(spec, crn.CheckTransactionAmount)
	if err != nil {
		return err
	}
	_, err = crn.c.AddFunc(spec, crn.CheckTransactionConfirmation)
	if err != nil {
		return err
	}
	return nil
}
