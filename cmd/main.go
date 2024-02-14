package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/everFinance/arseeding"
	"github.com/everFinance/arseeding/common"
	_ "github.com/mkevac/debugcharts"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "Arseeding",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "db_dir", Value: "./data/bolt", Usage: "bolt db dir path", EnvVars: []string{"DB_DIR"}},
			&cli.StringFlag{Name: "sqlite_dir", Value: "./data/sqlite", Usage: "sqlite db dir path", EnvVars: []string{"SQLITE_DIR"}},
			&cli.StringFlag{Name: "key_path", Value: "./data/bundler-keyfile.json", Usage: "ar keyfile path", EnvVars: []string{"KEY_PATH"}},
			&cli.StringFlag{Name: "ar_node", Value: "https://arweave.net", EnvVars: []string{"AR_NODE"}},
			&cli.StringFlag{Name: "pay", Value: "https://api-dev.everpay.io", Usage: "pay url", EnvVars: []string{"PAY"}},
			&cli.BoolFlag{Name: "no_fee", Value: false, EnvVars: []string{"NO_FEE"}},
			&cli.BoolFlag{Name: "manifest", Value: true, EnvVars: []string{"MANIFEST"}},
			&cli.IntFlag{Name: "bundle_interval", Value: 120, Usage: "bundle tx on chain time interval(seconds)", EnvVars: []string{"BUNDLE_INTERVAL"}},
			&cli.StringFlag{Name: "port", Value: ":8080", EnvVars: []string{"PORT"}},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	s := arseeding.New(
		c.String("db_dir"),
		c.String("sqlite_dir"),
		c.String("key_path"),
		c.String("ar_node"),
		c.String("pay"),
		c.Bool("no_fee"),
		c.Bool("manifest"),
		c.String("port"),
		c.Bool("use_kafka"), c.String("kafka_uri"))
	s.Run(c.String("port"), c.Int("bundle_interval"))

	common.NewMetricServer()

	<-signals

	s.Close()
	return nil
}
