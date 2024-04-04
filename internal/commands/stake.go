package commands

import (
	"fmt"
	"log"

	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
	"github.com/liteseed/edge/internal/contracts"
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

	ao, err := aogo.New()
	if err != nil {
		log.Fatal(err)
	}

	signer, err := goar.NewSignerFromPath(config.Signer)
	if err != nil {
		log.Fatal(err)
	}

	itemSigner, err := goar.NewItemSigner(signer)
	if err != nil {
		log.Fatal(err)
	}

	contract := contracts.New(ao, itemSigner)

	err = contract.Stake(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success")
	return nil
}
