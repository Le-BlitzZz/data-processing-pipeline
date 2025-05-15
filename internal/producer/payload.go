package producer

import (
	"Le-BlitzZz/streaming-etl-app/internal/producer/config"
	"encoding/json"
)

func createPayload(headers, row []string, split string) ([]byte, error) {
	p := make(map[string]any, len(headers)+1)
	p["split"] = split

	for i, header := range headers {
		p[header] = row[i]
	}

	return json.Marshal(p)
}

func publishPayload(conf *config.Config, payload []byte) error {
	if err := conf.Mb().Publish(conf.RawExchange(), payload); err != nil {
		return err
	}

	return conf.Mb().Publish(conf.ProcessedExchange(), payload)
}
