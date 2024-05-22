package commands

import (
	"log"
	"log/slog"
	"net/http"
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
	process := config.Process

	contracts := contracts.New(ao, process, itemSigner)

	c, err := cron.New(cron.WthContracts(contracts), cron.WithDatabase(db), cron.WithStore(store), cron.WithWallet(wallet), cron.WithLogger(logger))
	if err != nil {
		log.Fatal("failed to load cron", err)
	}
	err = c.Setup("* * * * *")
	if err != nil {
		log.Fatal("failed to setup cron", err)
	}

	go c.Start()

	s, err := server.New(":8080", context.App.Version, config.Node, server.WithContracts(contracts), server.WithDatabase(db), server.WithWallet(wallet), server.WithStore(store))
	if err != nil {
		log.Fatal("failed to setup server", err)
	}

	go func() {
		err := s.Start()
		if err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	<-quit

	log.Println("Shutdown")

	c.Shutdown()
	if err = db.Shutdown(); err != nil {
		log.Fatal("failed to shutdown database", err)
	}
	if err = store.Shutdown(); err != nil {
		log.Fatal("failed to shutdown store", err)
	}

	time.Sleep(2 * time.Second)
	if err = s.Shutdown(); err != nil {
		log.Fatal("failed to shutdown server", err)
	}
	return nil
}
