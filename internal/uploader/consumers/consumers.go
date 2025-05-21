package consumers

import (
	"Le-BlitzZz/streaming-etl-app/internal/uploader/entity"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/wagslane/go-rabbitmq"
)

var rawConsumer *rabbitmq.Consumer
var processedConsumer *rabbitmq.Consumer

func startHandlingWorkers(payloads chan []byte, factory entity.ApartmentFactory, onError func(error)) {
	for range runtime.NumCPU() {
		go func() {
			for payload := range payloads {
				if err := factory().CreateFromPayload(payload); err != nil {
					onError(err)
				}
			}
		}()
	}
}

func streamPayloadsFromConsumer(consumer *rabbitmq.Consumer, payloads chan<- []byte, onError func(error)) {
	if err := consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		payloads <- d.Body
		return rabbitmq.Ack
	}); err != nil {
		onError(err)
	}

	close(payloads)
}

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
