package commands

import (
	"log"

	"github.com/liteseed/bungo/api/routes"
	"github.com/liteseed/bungo/api/server"
	"github.com/liteseed/bungo/internal/database"
	"github.com/liteseed/bungo/internal/store"
	"github.com/urfave/cli/v2"
)

var Start = &cli.Command{
	Name:  "start",
	Usage: "Start the bundler node on this system",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "store", Value: "pebble", Usage: "store to use", EnvVars: []string{"STORE"}},
		&cli.StringFlag{Name: "database", Value: "sqlite", Usage: "database to use", EnvVars: []string{"DATABASE"}},
		&cli.StringFlag{Name: "key_path", Value: "./data/bundler-keyfile.json", Usage: "ar keyfile path", EnvVars: []string{"KEY_PATH"}},
		&cli.StringFlag{Name: "node", Value: "https://arweave.net", EnvVars: []string{"NODE"}},
		&cli.StringFlag{Name: "port", Value: ":8080", EnvVars: []string{"PORT"}},
	},
	Action: start,
}

func start(context *cli.Context) error {
	storeOption := context.String("STORE")
	sqlite := context.String("sqlite")

	database, err := database.New(sqlite, "sqlite")
	if err != nil {
		log.Fatal(err)
	}

	store := store.New(storeOption)

	a := routes.New(database, store)

	s := server.New()
	s.Register(a)
	s.Run(":8080")

	return nil
}
