package producer

import (
	"Le-BlitzZz/streaming-etl-app/internal/publisher/config"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) {
	start := time.Now()

	broker := conf.Broker()

	var err error

	rawPublisher, err = newPublisher(broker, conf.BrokerRawExchange())
	if err != nil {
		log.Panicf("failed to create publisher for raw messages: %s", err)
	}
	defer shutdownRawPublisher()

	processedPublisher, err = newPublisher(broker, conf.BrokerProcessedExchange())
	if err != nil {
		log.Panicf("failed to create publisher for processing raw messages: %s", err)
	}
	defer shutdownProcessedPublisher()

	go streamPayloadsFromCSVs(conf.DataSplitPathMap())

	runPublishers(conf.BrokerRawExchange(), conf.BrokerProcessedExchange())

	log.Infof("stream and publish work finished in [%s]", time.Since(start))
}
