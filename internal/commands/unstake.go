package commands

import (
	"log"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/sdk-go/contract"
	"github.com/urfave/cli/v2"
)

var Unstake = &cli.Command{
	Name:  "unstake",
	Usage: "Unstake the current bundler",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
	},
	Action: unstake,
}

func unstake(context *cli.Context) error {
	config := readConfig(context)

	process := config.Process
	s, err := signer.FromPath(config.Signer)
	if err != nil {
		return err
	}

	c := contract.New(process, s)

	res, err := c.Unstake()
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
