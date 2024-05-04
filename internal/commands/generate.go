package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"

	"github.com/everFinance/gojwk"
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
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatal(err)
	}
	jwk, err := gojwk.PrivateKey(key)
	if err != nil {
		log.Fatal(err)
	}
	data, err := gojwk.Marshal(jwk)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(config.Signer, data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
