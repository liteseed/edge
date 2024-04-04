package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	Balance,
	Migrate,
	Stake,
	Start,
	Unstake,
}

type Config struct {
	Port     string
	Process  string
	Signer   string
	Database string
	Store    string
	Node     string
}

func readConfig(context *cli.Context) Config {
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
	return config
}
