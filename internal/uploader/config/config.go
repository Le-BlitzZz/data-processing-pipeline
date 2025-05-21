package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

type Config struct {
	options *Options
	db      *gorm.DB
	broker  *rabbitmq.Conn
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	c := &Config{
		options: NewOptions(ctx),
	}

	if err := c.connectBroker(); err != nil {
		return nil, err
	}

	if err := c.connectDb(); err != nil {
		c.shutdownBroker()
		return nil, err
	}
	log.Info("config: successfully initialized")

	return c, nil
}

func (c *Config) Shutdown() {
	c.shutdownBroker()

	if err := c.closeDb(); err != nil {
		log.Errorf("could not close database connection: %s", err)
	} else {
		log.Info("closed database connection")
	}
}

func (c *Config) shutdownBroker() {
	if err := c.closeBroker(); err != nil {
		log.Errorf("could not close message broker connection: %s", err)
	} else {
		log.Info("closed message broker connection")
	}
}
