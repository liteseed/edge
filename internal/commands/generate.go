package commands

import (
	"os"

	"github.com/liteseed/goar/signer"
	"github.com/urfave/cli/v2"
)

var Generate = &cli.Command{
	Name:   "generate",
	Usage:  "Generate a new Arweave Private Wallet",
	Action: generate,
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
	},
}

func generate(ctx *cli.Context) error {
	config := readConfig(ctx)
	// Generate RSA key.
	s, err := signer.New()
	if err != nil {
		return err
	}

	err = os.WriteFile(config.Signer, s.PrivateKey.N.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
