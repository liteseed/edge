package commands

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
	"github.com/liteseed/edge/internal/contracts"
	"github.com/liteseed/edge/internal/cron"
	"github.com/liteseed/edge/internal/database"
	"github.com/liteseed/edge/internal/server"
	"github.com/liteseed/edge/internal/store"
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
	config := readConfig(ctx)

	logger := slog.New(slog.NewJSONHandler(&lumberjack.Logger{
		Filename:   config.Log,
		MaxSize:    2, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}, nil))

	db, err := database.New(config.Database)
	if err != nil {
		logger.Error(
			"failed: database connect",
			"error", err,
		)
		os.Exit(1)
	}

	wallet, err := goar.NewWalletFromPath(config.Signer, config.Node)
	if err != nil {
		logger.Error(
			"failed: wallet load",
			"error", err,
		)
		os.Exit(1)
	}
	itemSigner, err := goar.NewItemSigner(wallet.Signer)
	if err != nil {
		logger.Error(
			"failed: item-signer create",
			"error", err,
		)
		os.Exit(1)
	}

	store := store.New(config.Store)

	ao, err := aogo.New()
	if err != nil {
		logger.Error(
			"failed to connect to AO",
			"error", err,
		)
		os.Exit(1)
	}

	contracts := contracts.New(ao, itemSigner)

	c, err := cron.New(cron.WthContracts(contracts), cron.WithDatabase(db), cron.WithStore(store), cron.WithWallet(wallet), cron.WithLogger(logger))
	if err != nil {
		logger.Error(
			"failed: cron load",
			"error", err,
		)
		os.Exit(1)
	}
	err = c.Setup("* * * * *")
	if err != nil {
		logger.Error(
			"failed: cron load",
			"error", err,
		)
		os.Exit(1)
	}
	c.Start()

	s, err := server.New(":8080", server.WthContracts(contracts), server.WithDatabase(db), server.WithWallet(wallet.Signer), server.WithStore(store))
	if err != nil {
		logger.Error(
			"failed to start server",
			"error", err,
		)
		os.Exit(1)
	}
	go func() {
		if err = s.Start(); err != nil {
			logger.Error(
				"failed to start server",
				"error", err,
			)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown")

	c.Shutdown()
	if err = db.Shutdown(); err != nil {
		logger.Error(
			"failed to stop database",
			"error", err,
		)
		os.Exit(1)
	}
	if err = store.Shutdown(); err != nil {
		logger.Error(
			"failed to stop store",
			"error", err,
		)
		os.Exit(1)
	}

	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	time.Sleep(2 * time.Second)
	if err = s.Shutdown(cctx); err != nil {
		logger.Error(
			"failed to stop server",
			"error", err,
		)
		os.Exit(1)
	}
	return nil
}
