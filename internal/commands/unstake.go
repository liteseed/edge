package commands

import (
	"fmt"
	"log"

	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
	"github.com/liteseed/edge/internal/contracts"
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

	err = contract.Unstake()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success")
	return nil
}
