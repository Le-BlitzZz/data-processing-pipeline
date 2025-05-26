package uploader

import (
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/wagslane/go-rabbitmq"
)

var rawConsumer *rabbitmq.Consumer
var processedConsumer *rabbitmq.Consumer

func newConsumer(broker *rabbitmq.Conn, queue, exchange string) (*rabbitmq.Consumer, error) {
	return rabbitmq.NewConsumer(
		broker,
		queue,
		rabbitmq.WithConsumerOptionsExchangeName(exchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsRoutingKey(""),
		rabbitmq.WithConsumerOptionsExchangeKind("fanout"),
		rabbitmq.WithConsumerOptionsConcurrency(runtime.NumCPU()),
		rabbitmq.WithConsumerOptionsLogger(log.StandardLogger()),
		rabbitmq.WithConsumerOptionsExchangeDurable,
	)
}

func shutdownRawConsumer() {
	log.Info("Shutting down raw consumer")

	if rawConsumer != nil {
		rawConsumer.Close()
	}
}

func shutdownProcessedConsumer() {
	log.Info("Shutting down processed consumer")

	if processedConsumer != nil {
		processedConsumer.Close()
	}
}
