package commands

import (
	"fmt"
	"math"
	"strconv"

	"github.com/liteseed/goar/signer"

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

	s, err := signer.FromPath(config.Signer)
	if err != nil {
		return err
	}

	fmt.Println("Address: ", s.Address)

	c := contract.New(process, s)

	b, err := c.Balance(s.Address)
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

	pow := math.Pow10(denomination)

	bal, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		return err
	}

	_, err = fmt.Printf("Balance: %f %s\n", float64(bal)/pow, i.Ticker)
	if err != nil {
		return err
	}

	res, err := c.Staked()
	if err != nil {
		return err
	}

	_, err = fmt.Println("Staked: ", res)
	if err != nil {
		return err
	}

	return nil
}
