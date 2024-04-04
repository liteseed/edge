package commands

import (
	"log"

	"github.com/liteseed/edge/internal/database"
	"github.com/urfave/cli/v2"
)

var Migrate = &cli.Command{
	Name:  "migrate",
	Usage: "Run migration on your postgresql database",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
	},
	Action: migrate,
}

func migrate(context *cli.Context) error {
	config := readConfig(context)

	database, err := database.New(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migration Complete")
	return nil
}
