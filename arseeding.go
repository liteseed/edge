package arseeding

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/everFinance/arseeding/cache"
	"github.com/everFinance/arseeding/config"
	"github.com/everFinance/arseeding/schema"
	"github.com/everFinance/arseeding/sdk"
	"github.com/everFinance/go-everpay/common"
	paySdk "github.com/everFinance/go-everpay/sdk"
	"github.com/everFinance/goar"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

var log = common.NewLog("arseeding")

type Arseeding struct {
	store           *Store
	engine          *gin.Engine
	submitLocker    sync.Mutex
	endOffsetLocker sync.Mutex

	arCli     *goar.Client
	taskMg    *TaskManager
	scheduler *gocron.Scheduler

	cache    *Cache
	config   *config.Config
	KWriters map[string]*KWriter // key: topic

	// ANS-104 bundle
	arseedCli           *sdk.ArSeedCli
	everpaySdk          *paySdk.SDK
	wdb                 *Wdb
	bundler             *goar.Wallet
	bundlerItemSigner   *goar.ItemSigner
	NoFee               bool // if true, means no bundle fee; default false
	EnableManifest      bool
	bundlePerFeeMap     map[string]schema.Fee // key: tokenSymbol, val: fee per chunk_size(256KB)
	paymentExpiredRange int64                 // default
	expectedRange       int64                 // default 50 block
	locker              sync.RWMutex
	localCache          *cache.Cache
}

func New(
	boltDirPath, sqliteDir string,
	arWalletKeyPath string, arNode, payUrl string, noFee bool, enableManifest bool,
	port string, useKafka bool, kafkaUri string,
) *Arseeding {
	fmt.Println(boltDirPath)
	var err error

	KVDb, err := NewBoltStore(boltDirPath)

	if err != nil {
		panic(err)
	}

	jobmg := NewTaskMg()
	if err := jobmg.InitTaskMg(KVDb); err != nil {
		panic(err)
	}

	wdb := NewSqliteDb(sqliteDir)

	if err = wdb.Migrate(noFee, enableManifest); err != nil {
		panic(err)
	}
	bundler, err := goar.NewWalletFromPath(arWalletKeyPath, arNode)
	if err != nil {
		panic(err)
	}

	itemSigner, err := goar.NewItemSigner(bundler.Signer)
	if err != nil {
		panic(err)
	}
	everpaySdk, err := paySdk.New(bundler.Signer, payUrl)
	if err != nil {
		panic(err)
	}

	localArseedUrl := "http://127.0.0.1" + port
	a := &Arseeding{
		config:              config.New(sqliteDir),
		store:               KVDb,
		engine:              gin.Default(),
		submitLocker:        sync.Mutex{},
		endOffsetLocker:     sync.Mutex{},
		arCli:               goar.NewClient(arNode),
		taskMg:              jobmg,
		scheduler:           gocron.NewScheduler(time.UTC),
		arseedCli:           sdk.New(localArseedUrl),
		everpaySdk:          everpaySdk,
		wdb:                 wdb,
		bundler:             bundler,
		bundlerItemSigner:   itemSigner,
		NoFee:               noFee,
		EnableManifest:      enableManifest,
		bundlePerFeeMap:     make(map[string]schema.Fee),
		paymentExpiredRange: schema.DefaultPaymentExpiredRange,
		expectedRange:       schema.DefaultExpectedRange,
	}

	// init cache
	peerMap, err := KVDb.LoadPeers()
	if err != nil {
		peerMap = make(map[string]int64)
	}
	a.cache = NewCache(a.arCli, peerMap)
	if err := os.MkdirAll(schema.TmpFileDir, os.ModePerm); err != nil {
		panic(err)
	}

	if useKafka {
		kwriters, err := NewKWriters(kafkaUri)
		if err != nil {
			log.Error("NewKWriters(kafkaUri)", "err", err)
			panic(err)
		}
		a.KWriters = kwriters
	}

	localCache, err := cache.NewLocalCache(60 * time.Minute)
	if err != nil {
		log.Error("NewLocalCache", "err", err)
	}
	a.localCache = localCache
	return a
}

func (s *Arseeding) Run(port string, bundleInterval int) {
	s.config.Run()
	go s.runAPI(port)
	go s.runJobs(bundleInterval)
	go s.runTask()
}

func (s *Arseeding) Close() {
	s.store.Close()
	for _, k := range s.KWriters {
		k.Close()
	}
}

func (s *Arseeding) GetPerFee(tokenSymbol string) *schema.Fee {
	s.locker.RLock()
	defer s.locker.RUnlock()
	perFee, ok := s.bundlePerFeeMap[strings.ToUpper(tokenSymbol)]
	if !ok {
		return nil
	}
	return &perFee
}

func (s *Arseeding) SetPerFee(feeMap map[string]schema.Fee) {
	s.locker.Lock()
	s.bundlePerFeeMap = feeMap
	s.locker.Unlock()
}
