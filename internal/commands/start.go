package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/bungo/api/routes"
	"github.com/liteseed/bungo/api/server"
	"github.com/liteseed/bungo/internal/database"
	"github.com/liteseed/bungo/internal/store"
	"github.com/urfave/cli/v2"
)

type ConfigJSONValue struct {
	Name string 
	URL  string 
}

type Config struct {
	Port string
	Signer string
	Database ConfigJSONValue
	Store ConfigJSONValue
	Node ConfigJSONValue
}

var Start = &cli.Command{
	Name:  "start",
	Usage: "Start the bundler node on this system",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value", Required: true},
	},
	Action: start,
}

func start(context *cli.Context) error {
  configPath := context.Path("config")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal(err)
	}
	database, err := database.New(config.Database.Name, config.Database.URL)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := signer.New(config.Signer)
	if err != nil {
		log.Fatal(err)
	}

	store := store.New(config.Store.Name, config.Store.URL)
	a := routes.New(database, store, signer)

	s := server.New()
	s.Register(a)
	s.Run(":8080")

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
