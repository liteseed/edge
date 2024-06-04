package commands

import (
	"fmt"
	"math"
	"strconv"

	"github.com/everFinance/goar"

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

	process := config.Process

	signer, err := goar.NewSignerFromPath(config.Signer)
	if err != nil {
		return err
	}

	fmt.Println("Address: ", signer.Address)

	c := contract.New(process, signer)

	b, err := c.Balance(signer.Address)
	if err != nil {
		return err
	}

	i, err := c.Info()
	if err != nil {
		return err
	}
	denomination, err := strconv.Atoi(i.Denomination)
	if err != nil {
		return err
	}

	p := math.Pow10(denomination)

	bal, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		return err
	}

	res := float64(bal) / p

	_, err = fmt.Printf("Balance: %f %s\n", res, i.Ticker)
	if err != nil {
		return err
	}

	s, err := c.Staked()
	if err != nil {
		return err
	}

	_, err = fmt.Println("Staked: ", s)
	if err != nil {
		return err
	}

	return nil
}
