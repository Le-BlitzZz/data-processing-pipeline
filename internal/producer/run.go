package producer

import (
	"Le-BlitzZz/streaming-etl-app/internal/producer/config"
	"Le-BlitzZz/streaming-etl-app/pkg/csv"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) {
	start := time.Now()

	declareExchanges(conf)

	payloadCh := make(chan []byte, 256)
	var wg sync.WaitGroup

	startPublisherWorkers(conf, payloadCh, &wg)
	startCSVReaders(conf, payloadCh)

	wg.Wait()

	log.Infof("all datasets published to RabbitMQ after [%s]", time.Since(start))
}

func declareExchanges(conf *config.Config) {
	if err := conf.Mb().DeclareExchange(conf.RawExchange()); err != nil {
		log.Fatalf("failed to declare raw exchange: %s", err)
	}

	if err := conf.Mb().DeclareExchange(conf.ProcessedExchange()); err != nil {
		log.Fatalf("failed to declare processed exchange: %s", err)
	}
}

func startPublisherWorkers(conf *config.Config, payloadCh <-chan []byte, wg *sync.WaitGroup) {
	for range runtime.NumCPU() {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for payload := range payloadCh {
				if err := publishPayload(conf, payload); err != nil {
					log.Errorf("failed to publish payload: %s", err)
				}
			}
		}()
	}
}

func startCSVReaders(conf *config.Config, payloadCh chan<- []byte) {
	var wg sync.WaitGroup

	for _, split := range conf.Splits() {
		wg.Add(1)

		go func(split string) {
			defer wg.Done()

			rowHandler := rowProcessor(payloadCh, split)
			if err := csv.ForEachRow(conf.SplitPath(split), rowHandler, rowOnError); err != nil {
				log.Errorf("error processing split %s: %s", split, err)
			}
		}(split)
	}

	go func() {
		wg.Wait()
		close(payloadCh)

		log.Info("all splits sent to payload channel")
	}()
}
