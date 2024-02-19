package main

import (
	"log"
	"os"

	"github.com/liteseed/bungo/commands"
	"github.com/urfave/cli/v2"
)

const appName = "Bungo"
const appAbout = "Bungo"
const appEdition = "ce"
const appDescription = "bungo"

var version = "development"

// Metadata contains build specific information.
var Metadata = map[string]interface{}{
	"Name":        appName,
	"About":       appAbout,
	"Edition":     appEdition,
	"Description": appDescription,
	"Version":     version,
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()

	app := &cli.App{
		Name:        appName,
		Description: appDescription,
		Version:     version,
		Metadata:    Metadata,
		Commands:    commands.Bungo,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
