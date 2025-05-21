package publishers

import (
	"Le-BlitzZz/streaming-etl-app/pkg/csv"
	"encoding/json"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

func streamPayloadsFromCSV(splitPathMap map[string]string) {
	var count int64
	makeCountingHandler := func(split string) csv.HandlerFunc {
		return func(headers, row []string) {
			atomic.AddInt64(&count, 1)
			buildRowHandler(split)(headers, row)
		}
	}

	var wg sync.WaitGroup

	for split, path := range splitPathMap {
		wg.Add(1)

		go func(split, path string) {
			defer wg.Done()

			if err := csv.ForEachRow(path, makeCountingHandler(split), logRowError); err != nil {
				log.Errorf("error processing split %s: %s", split, err)
			}
		}(split, path)
	}

	wg.Wait()
	close(payloads)

	log.Infof("all splits processed: %d rows sent to payload channel", count)
}

func buildRowHandler(split string) csv.HandlerFunc {
	return func(headers, row []string) {
		payload, err := encodeRowToPayload(headers, row, split)
		if err != nil {
			log.Errorf("failed to create payload: %s", err)
			return
		}

		payloads <- payload
	}
}

func encodeRowToPayload(headers, row []string, split string) ([]byte, error) {
	p := make(map[string]any, len(headers)+1)
	p["split"] = split

	for i, header := range headers {
		p[header] = row[i]
	}

	return json.Marshal(p)
}

func logRowError(err error) {
	log.Errorf("failed to process row: %s", err)
}
