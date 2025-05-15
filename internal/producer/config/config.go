package config

import (
	"Le-BlitzZz/streaming-etl-app/internal/broker"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	mb *broker.MessageBroker
}

func NewConfig() (*Config, error) {
	c := &Config{}

	mbConfig := broker.NewConfig(
		c.BrokerUser(),
		c.BrokerPassword(),
		c.BrokerServer(),
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
