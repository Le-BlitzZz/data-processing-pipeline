package config

import (
	"Le-BlitzZz/streaming-etl-app/internal/broker"
	"log"
)

const (
	mbServer   = "rabbitmq:5672"
	mbUser     = "etlstream"
	mbPassword = "etlstream"
)

const (
	rawExchange       = "raw_exchange"
	processedExchange = "processed_exchange"
)

const datasetsDir = "datasets"

var splits = []string{"train", "test", "val"}

func (c *Config) Mb() *broker.MessageBroker {
	if c.mb == nil {
		log.Fatal("config: message broker not connected")
	}

	return c.mb
}

func (c *Config) RawExchange() string {
	return c.rawExchange
}

func (c *Config) ProcessedExchange() string {
	return c.processedExchange
}

func (c *Config) DatasetsDir() string {
	return c.datasetsDir
}

func (c *Config) Splits() []string {
	return c.splits
}
