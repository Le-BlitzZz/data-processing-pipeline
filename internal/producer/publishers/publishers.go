package publishers

import (
	"runtime"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/wagslane/go-rabbitmq"
)

var rawPublisher *rabbitmq.Publisher
var processedPublisher *rabbitmq.Publisher

func runPublishWorkers(rawExchange, processedExchange string) {
	var rawCount int64
	var processedCount int64

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

				if err := publishPayload(processedPublisher, payload, processedExchange); err != nil {
					log.Errorf("failed to publish payload: %s", err)
				} else {
					atomic.AddInt64(&processedCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	log.Infof("all payloads published: raw=%d processed=%d", rawCount, processedCount)
}

func newPublisher(broker *rabbitmq.Conn, exchange string) (*rabbitmq.Publisher, error) {
	return rabbitmq.NewPublisher(
		broker,
		rabbitmq.WithPublisherOptionsExchangeName(exchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind("fanout"),
		rabbitmq.WithPublisherOptionsLogger(log.StandardLogger()),
	)
}

func publishPayload(publisher *rabbitmq.Publisher, payload []byte, exchange string) error {
	return publisher.Publish(
		payload, []string{""},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exchange),
	)
}
