package commands

import (
	"log"

	"github.com/urfave/cli/v2"
)

var Benchmark = &cli.Command{
	Name:  "benchmark",
	Usage: "Run a benchmark on this system",
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
	Action: benchmark,
}

func benchmark(c *cli.Context) error {
	log.Fatal("not implemented")
	return nil
}
