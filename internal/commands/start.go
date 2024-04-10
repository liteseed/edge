package commands

import (
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
	}, &slog.HandlerOptions{AddSource: true}))

	db, err := database.New(config.Database)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	wallet, err := goar.NewWalletFromPath(config.Signer, config.Node)
	if err != nil {
		log.Fatal("failed to load wallet", "error", err)
	}
	itemSigner, err := goar.NewItemSigner(wallet.Signer)
	if err != nil {
		log.Fatal("failed to create item-signer", err)
	}

	store := store.New(config.Store)

	ao, err := aogo.New()
	if err != nil {
		log.Fatal("failed to connect to AO", err)
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

	time.Sleep(2 * time.Second)
	if err = s.Shutdown(); err != nil {
		logger.Error(
			"failed to stop server",
			"error", err,
		)
		os.Exit(1)
	}
	return nil
}
