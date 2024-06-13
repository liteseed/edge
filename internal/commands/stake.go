package commands

import (
	"log"

	"github.com/liteseed/goar/signer"

	"github.com/liteseed/sdk-go/contract"
	"github.com/urfave/cli/v2"
)

var Stake = &cli.Command{
	Name:  "stake",
	Usage: "Stake the current bundler",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
		&cli.StringFlag{Name: "url", Aliases: []string{"u"}, Usage: "url of bundler", Required: true},
	},
	Action: stake,
}

func stake(context *cli.Context) error {
	config := readConfig(context)
	url := context.String("url")

	s, err := signer.FromPath(config.Signer)
	if err != nil {
		return err
	}

	process := config.Process

	c := contract.New(process, s)

	res, err := c.Stake(url)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
