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

const appName = "Arseeding"
const appAbout = "appAbout"
const appEdition = "ce"
const appDescription = "Arseeding"

var version = "development"

// Metadata contains build specific information.
var Metadata = map[string]interface{}{
	"Name":        appName,
	"About":       appAbout,
	"Edition":     appEdition,
	"Description": appDescription,
	"Version":     version,
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()

	app := &cli.App{
		Name:        appName,
		Description: appDescription,
		Version:     version,
		Metadata:    Metadata,
		Commands: []*cli.Command{
			{
				Name: "deploy",
				Usage: "deploy the bundler",
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
			},
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
		c.String("bolt"),
		c.String("sqlite"),
		c.String("key_path"),
		c.String("node"),
		c.String("payment_url"),
		c.Bool("manifest"),
		c.String("port"),
		c.Bool("use_kafka"), c.String("kafka_uri"))
	s.Run(c.String("port"), c.Int("bundle_interval"))

	common.NewMetricServer()

	<-signals

	s.Close()
	return nil
}
