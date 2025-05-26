package main

import (
	"Le-BlitzZz/streaming-etl-app/internal/publisher/config"
	"Le-BlitzZz/streaming-etl-app/internal/publisher/producer"
	"os"

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
		Name:   "publisher",
		Usage:  "Start the Publisher",
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

	producer.Run(conf)

	return nil
}
