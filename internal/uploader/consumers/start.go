package consumers

import (
	"Le-BlitzZz/streaming-etl-app/internal/uploader/config"
	"Le-BlitzZz/streaming-etl-app/internal/uploader/entity"
	"context"

	log "github.com/sirupsen/logrus"
)

func Start(ctx context.Context, conf *config.Config) {
	consumer, err := newConsumer(conf.Broker(), conf.BrokerRawQueue(), conf.BrokerRawExchange())
	if err != nil {
		log.Fatalf("failed to create raw consumer: %s", err)
	}
	rawConsumer = consumer

	consumer, err = newConsumer(conf.Broker(), conf.BrokerProcessedQueue(), conf.BrokerProcessedExchange())
	if err != nil {
		log.Fatalf("failed to create processed consumer: %s", err)
	}
	processedConsumer = consumer

	go streamPayloadsFromConsumer(rawConsumer, rawPayloads, func(err error) {
		log.Errorf("failed to consume raw apartments: %s", err)
	})

	go streamPayloadsFromConsumer(processedConsumer, processedPayloads, func(err error) {
		log.Errorf("failed to consume processed apartments: %s", err)
	})

	startHandlingWorkers(rawPayloads, entity.NewRawApartment, func(err error) {
		log.Errorf("failed to handle raw apartment payload: %s", err)
	})
	startHandlingWorkers(processedPayloads, entity.NewProcessedApartment, func(err error) {
		log.Errorf("failed to handle processed apartment payload: %s", err)
	})

	<-ctx.Done()
	shutdown()
}

func shutdown() {
	log.Info("Shutting down raw consumer")

	if rawConsumer != nil {
		rawConsumer.Close()
	}

	log.Info("Shutting down processed consumer")

	if processedConsumer != nil {
		processedConsumer.Close()
	}
}
