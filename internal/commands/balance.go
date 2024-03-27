package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/liteseed/argo/ao"
	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
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

	var data = "Balance"
	var tags = []transaction.Tag{{Name: "Action", Value: "Balance"}}

	configPath := context.Path("config")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	var config Config

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalln(err)
	}

	signer, err := signer.New(config.Signer)
	if err != nil {
		log.Fatal(err)
	}

	ao := ao.New()

	messageId, err := ao.SendMessage(config.Process, data, tags, "", signer)
	if err != nil {
		log.Fatal(err)
	}

	result, err := ao.ReadResult(config.Process, messageId)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Printf("Balance: %s BUN\n", result.Messages[0]["Data"])
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
