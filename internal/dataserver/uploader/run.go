// uploader/run.go
package uploader

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"Le-BlitzZz/streaming-etl-app/internal/dataserver/config"
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/entity"

	"github.com/wagslane/go-rabbitmq"
)

const progressLogInterval = 30 * time.Second

func Run(ctx context.Context, conf *config.Config) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var err error
	broker := conf.Broker()

	rawConsumer, err = newConsumer(broker, conf.BrokerRawQueue(), conf.BrokerRawExchange())
	if err != nil {
		log.Panicf("failed to create raw consumer: %v", err)
	}
	processedConsumer, err = newConsumer(broker, conf.BrokerProcessedQueue(), conf.BrokerProcessedExchange())
	if err != nil {
		log.Panicf("failed to create processed consumer: %v", err)
		rawConsumer.Close()
		return
	}

	dataSize := conf.DataSize()

	var rawCount, processedCount int64
	increment := func(counter *int64) {
		if atomic.AddInt64(counter, 1) >= dataSize &&
			atomic.LoadInt64(&rawCount) >= dataSize &&
			atomic.LoadInt64(&processedCount) >= dataSize {
			cancel()
		}
	}

	rawHandler := func(d rabbitmq.Delivery) rabbitmq.Action {
		if err := (&entity.RawApartment{}).CreateFromPayload(d.Body); err != nil {
			log.Printf("failed to persist raw payload: %v", err)
			return rabbitmq.NackDiscard
		}
		increment(&rawCount)
		return rabbitmq.Ack
	}

	processedHandler := func(d rabbitmq.Delivery) rabbitmq.Action {
		if err := (&entity.ProcessedApartment{}).CreateFromPayload(d.Body); err != nil {
			log.Printf("failed to persist processed payload: %v", err)
			return rabbitmq.NackDiscard
		}
		increment(&processedCount)
		return rabbitmq.Ack
	}

	var progressWg sync.WaitGroup
	progressWg.Add(1)
	go func() {
		defer progressWg.Done()
		ticker := time.NewTicker(progressLogInterval)
		defer ticker.Stop()

		for {
			select {
			case <-cctx.Done():
				return
			case <-ticker.C:
				log.Printf("uploader: progress - RAW %d/%d | PROCESSED %d/%d",
					atomic.LoadInt64(&rawCount), dataSize,
					atomic.LoadInt64(&processedCount), dataSize)
			}
		}
	}()

	var consumersWg sync.WaitGroup
	consumersWg.Add(2)

	go func() {
		defer consumersWg.Done()
		if err := rawConsumer.Run(rawHandler); err != nil {
			log.Printf("raw consumer stopped with error: %v", err)
			cancel()
		}
	}()

	go func() {
		defer consumersWg.Done()
		if err := processedConsumer.Run(processedHandler); err != nil {
			log.Printf("processed consumer stopped with error: %v", err)
			cancel()
		}
	}()

	<-cctx.Done()
	shutdownRawConsumer()
	shutdownProcessedConsumer()
	conf.ShutdownBroker()
	consumersWg.Wait()
	progressWg.Wait()

	log.Printf("uploader: finished - RAW %d/%d | PROCESSED %d/%d", rawCount, dataSize, processedCount, dataSize)
}
