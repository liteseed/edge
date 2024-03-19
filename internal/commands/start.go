package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/edge/api/routes"
	"github.com/liteseed/edge/api/server"
	"github.com/liteseed/edge/internal/cron"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/store"
	"github.com/urfave/cli/v2"
)

type JSONValue struct {
	Name string
	URL  string
}

type StartConfig struct {
	Port     string
	Signer   string
	Database JSONValue
	Store    JSONValue
	Node     JSONValue
}

var Start = &cli.Command{
	Name:  "start",
	Usage: "Start the bundler node on this system",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
	},
	Action: start,
}

func start(context *cli.Context) error {
	configPath := context.Path("config")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config StartConfig

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

	c, err := cron.New(cron.WithDatabase(database), cron.WithSigner(signer), cron.WithStore(store))
	if err != nil {
		log.Fatalln("failed to load cron", err)
	}
	err = c.Add("* * * * *")
	if err != nil {
		log.Fatalln("failed to load cron", err)
	}
	c.Start()

	s := server.New()
	s.Register(a)
	s.Run(":8080")

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
