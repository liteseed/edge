package commands

import (
	"fmt"
	"log"

	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
	"github.com/liteseed/edge/internal/contracts"
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

	signer, err := goar.NewSignerFromPath(config.Signer)
	if err != nil {
		log.Fatal(err)
	}

	itemSigner, err := goar.NewItemSigner(signer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Address: ", signer.Address)

	contract := contracts.New(ao, itemSigner)

	b, err := contract.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Printf("Balance: %s BUN\n", b)
	if err != nil {
		log.Fatal(err)
	}

	s, err := contract.GetStaker()
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Println("Staked: ", s)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
