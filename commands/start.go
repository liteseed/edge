package commands

import (
	"log"
	"time"

	"github.com/liteseed/bungo/cache"
	"github.com/liteseed/bungo/database"
	"github.com/liteseed/bungo/server"
	"github.com/liteseed/bungo/store"
	"github.com/urfave/cli/v2"
)

var Start = &cli.Command{
	Name:  "start",
	Usage: "Start the bundler on this system",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "bolt", Value: "./data/bolt", Usage: "bolt db dir path", EnvVars: []string{"BOLT"}},
		&cli.StringFlag{Name: "sqlite", Value: "./data/sqlite", Usage: "sqlite db dir path", EnvVars: []string{"SQLITE"}},
		&cli.StringFlag{Name: "key_path", Value: "./data/bundler-keyfile.json", Usage: "ar keyfile path", EnvVars: []string{"KEY_PATH"}},
		&cli.StringFlag{Name: "node", Value: "https://arweave.net", EnvVars: []string{"NODE"}},
		&cli.StringFlag{Name: "payment_url", Value: "https://api-dev.everpay.io", Usage: "pay url", EnvVars: []string{"PAYMENT_URL"}},
		&cli.BoolFlag{Name: "manifest", Value: true, EnvVars: []string{"MANIFEST"}},
		&cli.IntFlag{Name: "bundle_interval", Value: 120, Usage: "bundle tx on chain time interval(seconds)", EnvVars: []string{"BUNDLE_INTERVAL"}},
		&cli.StringFlag{Name: "port", Value: ":8080", EnvVars: []string{"PORT"}},
	},
	Action: start,
}

func start(context *cli.Context) error {
	bolt := context.String("bolt")
	sqlite := context.String("sqlite")

	cache, err := cache.NewBigCache(60 * time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewSqliteDatabase(sqlite)
	if err = db.Migrate(); err != nil {
		log.Fatal(err)
	}
	store, err := store.NewBoltStore(bolt)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(cache, db, store)
	s.Run(":8080")
	return nil
}
