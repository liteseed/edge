package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/everFinance/arseeding"
	"github.com/everFinance/arseeding/common"
	"github.com/urfave/cli/v2"
)

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
