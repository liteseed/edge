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
}

func generate(ctx *cli.Context) error {
	filename := "signer.json"
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
	err = os.WriteFile(filename, data, os.ModePerm)
	return nil
}
