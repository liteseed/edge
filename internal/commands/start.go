package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/cron"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/server"
	"github.com/liteseed/edge/internal/store"
	"github.com/urfave/cli/v2"
)

type JSONValue struct {
	Name string
	URL  string
}

type Config struct {
	Port     string
	Process  string
	Signer   string
	Database string
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
		log.Fatalln(err)
	}

	var config Config

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalln(err)
	}

	database, err := database.New(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := goar.NewWalletFromPath(config.Signer, config.Node.URL)
	if err != nil {
		log.Fatal(err)
	}

	store := store.New(config.Store.Name, config.Store.URL)

	ao, err := aogo.New()
	if err != nil {
		log.Fatalln("failed to load ao", err)
	}
	itemSigner, err := goar.NewItemSigner(wallet.Signer)
	if err != nil {
		log.Fatalln("failed to load ao", err)
	}
	
	contracts := contracts.New(ao, itemSigner)

	c, err := cron.New(cron.WthContracts(contracts), cron.WithDatabase(database), cron.WithWallet(wallet), cron.WithStore(store))
	if err != nil {
		log.Fatalln("failed to load cron", err)
	}
	err = c.PostBundle("* * * * *")
	if err != nil {
		log.Fatalln("failed to start bundle posting service", err)
	}
	err = c.Notify()
	if err != nil {
		log.Fatalln("failed to start notification service", err)
	}
	c.Start()

	s := server.New(contracts, database, store)
	s.Run(":8080")

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
