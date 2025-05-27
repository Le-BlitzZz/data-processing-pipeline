package producer

import (
	"runtime"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/wagslane/go-rabbitmq"
)

var rawPublisher *rabbitmq.Publisher
var rawProcessingPublisher *rabbitmq.Publisher

func runPublishers(rawExchange, rawProcessingExchange string) {
	var rawCount int64
	var rawProcessingCount int64

	var wg sync.WaitGroup

	for range runtime.NumCPU() {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for payload := range payloads {
				if err := publishPayload(rawPublisher, payload, rawExchange); err != nil {
					log.Errorf("failed to publish payload: %s", err)
				} else {
					atomic.AddInt64(&rawCount, 1)
				}

				if err := publishPayload(rawProcessingPublisher, payload, rawProcessingExchange); err != nil {
					log.Errorf("failed to publish payload: %s", err)
				} else {
					atomic.AddInt64(&rawProcessingCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	log.Infof("all payloads published: raw=%d rawProcessing=%d", rawCount, rawProcessingCount)
}

func newPublisher(broker *rabbitmq.Conn, exchange string) (*rabbitmq.Publisher, error) {
	return rabbitmq.NewPublisher(
		broker,
		rabbitmq.WithPublisherOptionsExchangeName(exchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind("fanout"),
		rabbitmq.WithPublisherOptionsLogger(log.StandardLogger()),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
}

func publishPayload(publisher *rabbitmq.Publisher, payload []byte, exchange string) error {
	return publisher.Publish(
		payload, []string{""},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exchange),
	)
}

func shutdownRawPublisher() {
	log.Info("shutting down raw publisher")

	if rawPublisher != nil {
		rawPublisher.Close()
	}
}

func shutdownProcessedPublisher() {
	log.Info("shutting down processed publisher")

	if rawProcessingPublisher != nil {
		rawProcessingPublisher.Close()
	}
}
