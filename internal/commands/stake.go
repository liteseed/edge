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

var Stake = &cli.Command{
	Name:  "stake",
	Usage: "Stake the current bundler",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
		&cli.PathFlag{Name: "url", Value: "https://edge.liteseed.xyz", Usage: "url of bundler"},
	},
	Action: stake,
}

func stake(context *cli.Context) error {

	var data = "Stake"
	var tags = []types.Tag{{Name: "Action", Value: "Stake"}}

	configPath := context.Path("config")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	url := context.Path("url")

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
	tags = append(tags, types.Tag{Name: "URL", Value: url})
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
