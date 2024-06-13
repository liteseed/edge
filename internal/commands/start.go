package commands

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liteseed/edge/internal/cron"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/server"
	"github.com/liteseed/edge/internal/store"
	"github.com/liteseed/goar/client"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/sdk-go/contract"
	"github.com/urfave/cli/v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Start = &cli.Command{
	Name:  "start",
	Usage: "Start the bundler node on this system",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.json", Usage: "path to config value"},
	},
	Action: start,
}

func start(context *cli.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config := readConfig(context)

	logger := slog.New(
		slog.NewJSONHandler(
			&lumberjack.Logger{
				Filename:   config.Log,
				MaxSize:    2,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
			},
			&slog.HandlerOptions{AddSource: true},
		),
	)

	db, err := database.New(config.Database)
	if err != nil {
		return err
	}

	signer, err := signer.FromPath(config.Signer)
	if err != nil {
		return err
	}

	client := client.New(config.Node)

	store := store.New(config.Store)

	process := config.Process

	contracts := contract.New(process, signer)

	cron, err := cron.New(cron.WithClient(client), cron.WithContracts(contracts), cron.WithDatabase(db), cron.WithStore(store), cron.WithSigner(signer), cron.WithLogger(logger))
	if err != nil {
		return err
	}
	err = cron.Setup("* * * * *")
	if err != nil {
		return err
	}

	go cron.Start()

	server, err := server.New(":8080", context.App.Version, server.WithClient(client), server.WithContracts(contracts), server.WithDatabase(db), server.WithSigner(signer), server.WithStore(store))
	if err != nil {
		return err
	}

	go func() {
		err := server.Start()
		if err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	<-quit

	log.Println("Shutdown")

	cron.Shutdown()
	if err = db.Shutdown(); err != nil {
		return err
	}
	if err = store.Shutdown(); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	if err = server.Shutdown(); err != nil {
		return err
	}
	return nil
}
