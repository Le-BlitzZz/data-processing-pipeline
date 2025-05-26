package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wagslane/go-rabbitmq"
)

type Config struct {
	options *Options
	broker  *rabbitmq.Conn
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	c := &Config{
		options: NewOptions(ctx),
	}

	if err := c.connectBroker(); err != nil {
		return nil, err
	}

	log.Info("config: successfully initialized")

	return c, nil
}

func (c *Config) Shutdown() {
	if err := c.closeBroker(); err != nil {
		log.Errorf("could not close message broker connection: %s", err)
	} else {
		log.Info("closed message broker connection")
	}
}
