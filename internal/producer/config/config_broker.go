package config

import (
	"Le-BlitzZz/streaming-etl-app/internal/broker"
	"log"
)

const (
	mbServer   = "rabbitmq:5672"
	mbUser     = "etlstream"
	mbPassword = "etlstream"

	mbRawExchange       = "raw_exchange"
	mbProcessedExchange = "processed_exchange"
)

func (c *Config) Mb() *broker.MessageBroker {
	if c.mb == nil {
		log.Fatal("config: message broker not connected")
	}

	return c.mb
}

func (c *Config) BrokerUser() string {
	return mbUser
}

func (c *Config) BrokerPassword() string {
	return mbPassword
}

func (c *Config) BrokerServer() string {
	return mbServer
}

func (c *Config) RawExchange() string {
	return mbRawExchange
}

func (c *Config) ProcessedExchange() string {
	return mbProcessedExchange
}
