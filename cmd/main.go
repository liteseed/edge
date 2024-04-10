package main

import (
	"log"
	"os"

	"github.com/liteseed/edge/internal/commands"
	"github.com/urfave/cli/v2"
)

const Name = "Edge"
const Description = "Edge is the bundler node for the Liteseed Network"

var Version string

func main() {

	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()

	app := &cli.App{
		Name:           Name,
		Usage:          "Go to https://docs.liteseed.xyz to get started",
		Description:    Description,
		Version:        Version,
		Commands:       commands.Commands,
		DefaultCommand: "help",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
