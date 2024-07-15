package commands

import (
	"github.com/liteseed/edge/internal/database"
	"github.com/urfave/cli/v2"
	"log"
)

var Migrate = &cli.Command{
	Name:  "migrate",
	Usage: "Run migration on your postgresql database",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
	},
	Action: migrate,
}

func migrate(ctx *cli.Context) error {
	config := readConfig(ctx)

	db, err := database.New(config.Driver, config.Database)
	if err != nil {
		return err
	}

	err = db.Migrate()
	if err != nil {
		return err
	}

	log.Println("Migration Complete")
	err = db.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
