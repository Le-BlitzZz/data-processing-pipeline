package main

import (
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/config"
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/get"
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/presenter"
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/uploader"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()

	app := &cli.App{
		Name:   "dataserver",
		Usage:  "Start the Dataserver",
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func run(ctx *cli.Context) error {
	conf, err := initConfig(ctx)
	if err != nil {
		return err
	}
	defer conf.Shutdown()

	cctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTSTP, syscall.SIGTERM)
	defer cancel()

	conf.InitDb()

	go presenter.Start(cctx, conf)

	go uploader.Run(cctx, conf)

	<-cctx.Done()

	log.Info("shutting down...")
	cancel()

	time.Sleep(1 * time.Second)

	return nil
}

func initConfig(ctx *cli.Context) (*config.Config, error) {
	conf, err := config.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	get.SetConfig(conf)
	return conf, nil
}
