package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/liteseed/argo/ao"
	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
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

	var data = "Unstake"
	var tags = []transaction.Tag{{Name: "Action", Value: "Unstake"}}

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

	_, err = ao.ReadResult(config.Process, messageId)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Success")
	return nil
}
