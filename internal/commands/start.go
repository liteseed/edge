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
	"github.com/liteseed/goar/wallet"
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

func start(ctx *cli.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config := readConfig(ctx)

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

	w, err := wallet.FromPath(config.Signer, config.Node)
	if err != nil {
		return err
	}

	store := store.New(config.Store)

	process := config.Process

	contracts := contract.New(process, w.Signer)

	crn, err := cron.New(cron.WithContracts(contracts), cron.WithDatabase(db), cron.WithStore(store), cron.WithWallet(w), cron.WithLogger(logger))
	if err != nil {
		return err
	}
	err = crn.Setup("* * * * *")
	if err != nil {
		return err
	}

	go crn.Start()

	srv, err := server.New(":8080", ctx.App.Version, server.WithContracts(contracts), server.WithDatabase(db), server.WithWallet(w), server.WithStore(store))
	if err != nil {
		return err
	}

	go func() {
		err := srv.Start()
		if err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	<-quit

	log.Println("Shutdown")

	crn.Shutdown()
	if err = db.Shutdown(); err != nil {
		return err
	}
	if err = store.Shutdown(); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	if err = srv.Shutdown(); err != nil {
		return err
	}
	return nil
}
