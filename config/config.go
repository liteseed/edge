package config

import (
	"time"

	"github.com/everFinance/go-everpay/account"
	"github.com/go-co-op/gocron"
	"github.com/liteseed/bungo/config/schema"
	"github.com/liteseed/bungo/store"
)

type Config struct {
	store          *store.Store
	wdb            *Wdb
	speedTxFee     int64
	bundleServeFee int64
	ipWhiteList    map[string]struct{}
	scheduler      *gocron.Scheduler
	Param          schema.Param
}

func New(bolt string, sqlite string) *Config {
	store, err := store.NewBoltStore(bolt)
	if err != nil {
		panic(err)
	}
	wdb := NewSqliteDb(sqlite)
	err = wdb.Migrate()
	if err != nil {
		panic(err)
	}
	fee, err := wdb.GetFee()
	if err != nil {
		panic(err)
	}
	param, err := wdb.GetParam()
	if err != nil {
		panic(err)
	}
	return &Config{
		store:          store,
		wdb:            wdb,
		speedTxFee:     fee.SpeedTxFee,
		bundleServeFee: fee.BundleServeFee,
		ipWhiteList:    make(map[string]struct{}),
		scheduler:      gocron.NewScheduler(time.UTC),
		Param:          param,
	}
}

func (c *Config) GetSpeedFee() int64 {
	return c.speedTxFee
}

func (c *Config) GetServeFee() int64 {
	return c.bundleServeFee
}

func (c *Config) GetIPWhiteList() *map[string]struct{} {
	return &c.ipWhiteList
}

func (c *Config) Run() {
	go c.runJobs()
}

func (c *Config) Close() {
	c.wdb.Close()
}

func (s *Config) FeeCollectAddress() string {
	feeCfg, err := s.wdb.GetFee()
	if err != nil {
		return ""
	}
	collectAddr := feeCfg.FeeCollectAddress
	_, accId, err := account.IDCheck(collectAddr)
	if err != nil {
		log.Error("fee collection address incorrect", "err", err, "addr", collectAddr)
		return ""
	}
	return accId
}
