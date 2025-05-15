package main

import (
	"Le-BlitzZz/streaming-etl-app/internal/entity"
	"Le-BlitzZz/streaming-etl-app/internal/uploader/config"
	"encoding/json"
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
		Name:   "uploader",
		Usage:  "Start the Uploader",
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func run(cli *cli.Context) error {
	conf, err := config.NewConfig()
	if err != nil {
		return err
	}
	defer conf.Shutdown()

	conf.Db().Init()

	mb := conf.Mb()

	if err := mb.DeclareExchange(conf.RawExchange()); err != nil {
		return err
	}

	// if err := mb.DeclareExchange(conf.ProcessedExchange()); err != nil {
	// 	return err
	// }

	q, err := mb.DeclareQueue(conf.RawQueue())
	if err != nil {
		return err
	}

	if err := mb.BindQueue(q.Name, conf.RawExchange()); err != nil {
		return err
	}

	deliveries, err := mb.Consume(q.Name)
	if err != nil {
		return err
	}
	for d := range deliveries {
		if err := processRawPayload(conf, d.Body); err != nil {
			d.Nack(false, true)
		} else {
			d.Ack(true)
		}
	}

	return nil
}

func processRawPayload(conf *config.Config, data []byte) error {
	rawApartment := &entity.RawApartment{}
	if err := json.Unmarshal(data, rawApartment); err != nil {
		return err
	}

	return conf.Db().Create(rawApartment)
}

// func processProcessedPayload(conf *config.Config, data []byte) error {
// 	processedApartment := &entity.ProcessedApartment{}
// 	if err := json.Unmarshal(data, processedApartment); err != nil {
// 		return err
// 	}

// 	return conf.Db().Create(processedApartment)
// }
