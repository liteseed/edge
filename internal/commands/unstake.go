package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/liteseed/aogo"
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
	var tags = []types.Tag{{Name: "Action", Value: "Unstake"}}

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

	messageId, err := ao.SendMessage(config.Process, data, tags, "", itemSigner)
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
