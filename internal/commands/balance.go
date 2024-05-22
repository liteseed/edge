package commands

import (
	"fmt"
	"log"

	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"

	"github.com/liteseed/sdk-go/contract"
	"github.com/urfave/cli/v2"
)

var Balance = &cli.Command{
	Name:  "balance",
	Usage: "Check the balance of the wallet",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
	},
	Action: balance,
}

func balance(context *cli.Context) error {
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

	fmt.Println("Address: ", signer.Address)

	c := contract.New(ao, process, itemSigner)

	b, err := c.Balance(signer.Address)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Printf("Balance: %s LSD\n", b)
	if err != nil {
		log.Fatal(err)
	}

	s, err := c.Staked()
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Println("Staked: ", s)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
