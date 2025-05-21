package publishers

import (
	"Le-BlitzZz/streaming-etl-app/internal/producer/config"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) {
	start := time.Now()

	publisher, err := newPublisher(conf.Broker(), conf.BrokerRawExchange())
	if err != nil {
		log.Fatalf("failed to create publisher for raw messages: %s", err)
	}
	rawPublisher = publisher

	publisher, err = newPublisher(conf.Broker(), conf.BrokerProcessedExchange())
	if err != nil {
		log.Fatalf("failed to create publisher for processing raw messages: %s", err)
	}
	processedPublisher = publisher

	go streamPayloadsFromCSV(conf.SplitPathMap())

	runPublishWorkers(conf.BrokerRawExchange(), conf.BrokerProcessedExchange())

	shutdown()

	log.Infof("stream and publish work finished in [%s]", time.Since(start))
}

func shutdown() {
	log.Info("shutting down raw publisher")

	if rawPublisher != nil {
		rawPublisher.Close()
	}

	log.Info("shutting down processed publisher")

	if processedPublisher != nil {
		processedPublisher.Close()
	}
}
