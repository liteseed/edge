package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	Balance,
	Generate,
	Migrate,
	Stake,
	Start,
	Unstake,
}

type Config struct {
	Database string
	Log      string
	Port     string
	Process  string
	Node     string
	Signer   string
	Store    string
}

func readConfig(ctx *cli.Context) Config {
	configPath := ctx.Path("config")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	var config Config

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
