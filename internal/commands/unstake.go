package commands

import (
	"fmt"
	"log"

	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
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
	ao, err := aogo.New()
	if err != nil {
		log.Fatal(err)
	}
	process := config.Process
	signer, err := goar.NewSignerFromPath(config.Signer)
	if err != nil {
		log.Fatal(err)
	}
	itemSigner, err := goar.NewItemSigner(signer)
	if err != nil {
		log.Fatal(err)
	}
	c := contract.New(ao, process, itemSigner)

	err = c.Unstake()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success")
	return nil
}
