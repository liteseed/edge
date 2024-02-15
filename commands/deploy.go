package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/liteseed/bungo"
	"github.com/urfave/cli/v2"
)

var Deploy = &cli.Command{
	Name:  "deploy",
	Usage: "Deploy the bundler on this system",
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
	Action: deploy,
}

func deploy(context *cli.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	server := bungo.New(
		context.String("bolt"),
		context.String("sqlite"),
		context.String("key_path"),
		context.String("node"),
		context.String("payment_url"),
		context.Bool("manifest"),
		context.String("port"),
		context.Bool("use_kafka"), context.String("kafka_uri"))

	server.Run(context.String("port"), context.Int("bundle_interval"))

	<-signals

	server.Close()
	return nil
}
