package producer

import (
	"Le-BlitzZz/streaming-etl-app/pkg/csv"

	log "github.com/sirupsen/logrus"
)

func rowProcessor(payloadCh chan<- []byte, split string) csv.HandlerFunc {
	return func(headers, row []string) {
		payload, err := createPayload(headers, row, split)
		if err != nil {
			log.Errorf("failed to create payload: %s", err)
			return
		}

		payloadCh <- payload
	}
}

func rowOnError(err error) {
	log.Errorf("failed to process row: %s", err)
}
