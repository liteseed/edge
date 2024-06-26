package cron

import (
	"log/slog"

	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/liteseed/goar/wallet"
	"github.com/liteseed/sdk-go/contract"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	cron     *cron.Cron
	contract *contract.Contract
	database *database.Database
	logger   *slog.Logger
	wallet   *wallet.Wallet
	store    *store.Store
}

type Option = func(*Cron)

func New(options ...func(*Cron)) (*Cron, error) {
	c := &Cron{cron: cron.New()}
	for _, o := range options {
		o(c)
	}
	return c, nil
}

func WithContracts(contract *contract.Contract) Option {
	return func(crn *Cron) {
		crn.contract = contract
	}
}

func WithDatabase(db *database.Database) Option {
	return func(crn *Cron) {
		crn.database = db
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(crn *Cron) {
		crn.logger = logger
	}
}

func WithWallet(w *wallet.Wallet) Option {
	return func(crn *Cron) {
		crn.wallet = w
	}
}

func WithStore(s *store.Store) Option {
	return func(crn *Cron) {
		crn.store = s
	}
}

func (crn *Cron) Start() {
	crn.cron.Start()
}

func (crn *Cron) Shutdown() {
	crn.cron.Stop()
}

func (crn *Cron) Setup(spec string) error {
	_, err := crn.cron.AddFunc(spec, crn.JobBundleConfirmations)
	if err != nil {
		return err
	}
	_, err = crn.cron.AddFunc(spec, crn.JobPostBundle)
	if err != nil {
		return err
	}
	_, err = crn.cron.AddFunc(spec, crn.JobPostUpdate)
	if err != nil {
		return err
	}
	_, err = crn.cron.AddFunc(spec, crn.JobReleaseReward)
	if err != nil {
		return err
	}
	_, err = crn.cron.AddFunc(spec, crn.JobDeleteBundle)
	if err != nil {
		return err
	}
	return nil
}
