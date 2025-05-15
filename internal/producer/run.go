package producer

import (
	"Le-BlitzZz/streaming-etl-app/internal/producer/config"
	"Le-BlitzZz/streaming-etl-app/pkg/csv"
	"encoding/json"
	"path/filepath"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
)

func Run(conf *config.Config) error {
	mb := conf.Mb()

	if err := mb.DeclareExchange(conf.RawExchange()); err != nil {
		return err
	}

	// if err := mb.DeclareExchange(config.ProcessedExchange); err != nil {
	// 	return err
	// }

	rowsCh := make(chan []byte, 1000)

	var wg sync.WaitGroup

	startPublisherWorkers(&wg, conf, rowsCh, runtime.NumCPU())

	streamSplits(conf, rowsCh)

	close(rowsCh)
	wg.Wait()

	log.Info("all rows published")
	return nil
}

func startPublisherWorkers(wg *sync.WaitGroup, conf *config.Config, rows <-chan []byte, count int) {
	for i := range count {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for payload := range rows {
				if err := conf.Mb().Publish(conf.RawExchange(), payload); err != nil {
					log.Errorf("failed to publish to exchange: %s", err)
				}

				// if err := conf.Mb().Publish(conf.ProcessedExchange, payload); err != nil {
				// 	log.Errorf("failed to publish to exchange: %s", err)
				// }
			}
		}(i)
	}
}

func streamSplits(conf *config.Config, rowsCh chan<- []byte) {
	var wg sync.WaitGroup

	for _, split := range conf.Splits() {
		wg.Add(1)
		go func(split string) {
			defer wg.Done()

			csv.ForEachRow(getSplitFile(conf.DatasetsDir(), split), func(headers, row []string) {
				payload, err := createPayload(headers, row, split)
				if err != nil {
					log.Errorf("failed to create payload: %s", err)
					return
				}

				rowsCh <- payload
			})
		}(split)
	}

	wg.Wait()
}

func createPayload(headers, row []string, split string) ([]byte, error) {
	p := make(map[string]any, len(headers)+1)
	p["split"] = split

	for i, header := range headers {
		p[header] = row[i]
	}

	return json.Marshal(p)
}

func getSplitFile(datasetDir, split string) string {
	return filepath.Join(datasetDir, split+".csv")
}
