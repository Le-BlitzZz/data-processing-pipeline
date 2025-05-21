package main

import (
	"Le-BlitzZz/streaming-etl-app/internal/uploader/config"
	"Le-BlitzZz/streaming-etl-app/internal/uploader/consumers"
	"Le-BlitzZz/streaming-etl-app/internal/uploader/server"
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
		Name:   "uploader",
		Usage:  "Start the Uploader",
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func run(ctx *cli.Context) error {
	conf, err := config.NewConfig(ctx)
	if err != nil {
		return err
	}
	defer conf.Shutdown()

	cctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTSTP)

	conf.InitDb()

	go server.Start(cctx, conf)

	go consumers.Start(cctx, conf)

	<-cctx.Done()

	log.Info("shutting down...")
	cancel()

	time.Sleep(2 * time.Second)
	conf.Shutdown()

	return nil
}
