package config

import (
	"Le-BlitzZz/streaming-etl-app/internal/broker"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	mb *broker.MessageBroker

	rawExchange       string
	processedExchange string
	datasetsDir       string
	splits            []string
}

func NewConfig() (*Config, error) {
	c := &Config{
		rawExchange:       rawExchange,
		processedExchange: processedExchange,
		datasetsDir:       datasetsDir,
		splits:            splits,
	}

	mbConfig := broker.NewConfig(
		mbUser,
		mbPassword,
		mbServer,
	)

	mb, err := broker.NewMessageBroker(mbConfig)
	if err != nil {
		return nil, err
	}

	c.mb = mb

	log.Debug("config: successfully initialized")

	return c, nil
}

func (c *Config) Shutdown() {
	if err := c.mb.Close(); err != nil {
		log.Errorf("could not close message broker connection: %s", err)
	} else {
		log.Debug("closed message broker connection")
	}
}
