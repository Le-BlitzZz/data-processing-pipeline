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

const (
	rawQueue       = "raw_queue"
	processedQueue = "processed_queue"
)

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

func (c *Config) RawQueue() string {
	return c.rawQueue
}

func (c *Config) ProcessedQueue() string {
	return c.processedQueue
}
