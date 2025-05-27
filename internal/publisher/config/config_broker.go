package config

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
)

const (
	brokerServer   = "rabbitmq:5672"
	brokerUser     = "etlstream"
	brokerPassword = "etlstream"

	brokerRawExchange           = "raw_exchange"
	brokerRawProcessingExchange = "raw_processing_exchange"
)

func (c *Config) Broker() *rabbitmq.Conn {
	if c.broker == nil {
		log.Fatal("config: message broker not connected")
	}

	return c.broker
}

func (c *Config) BrokerUser() string {
	if c.options.BrokerUser == "" {
		return brokerUser
	}
	return c.options.BrokerUser
}

func (c *Config) BrokerPassword() string {
	if c.options.BrokerPassword == "" {
		return brokerPassword
	}
	return c.options.BrokerPassword
}

func (c *Config) BrokerServer() string {
	if c.options.BrokerServer == "" {
		return brokerServer
	}
	return c.options.BrokerServer
}

func (c *Config) BrokerRawExchange() string {
	if c.options.BrokerRawExchange == "" {
		return brokerRawExchange
	}
	return c.options.BrokerRawExchange
}

func (c *Config) BrokerRawProcessingExchange() string {
	if c.options.BrokerRawProcessingExchange == "" {
		return brokerRawProcessingExchange
	}
	return c.options.BrokerRawProcessingExchange
}

func (c *Config) BrokerDsn() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s/",
		c.BrokerUser(),
		c.BrokerPassword(),
		c.BrokerServer(),
	)
}

func (c *Config) connectBroker() error {
	brokerDsn := c.BrokerDsn()

	broker, err := rabbitmq.NewConn(brokerDsn, rabbitmq.WithConnectionOptionsLogger(log.StandardLogger()))
	if err != nil || broker == nil {
		log.Infof("config: waiting for the message broker to become available")

		for range 5 {
			broker, err = rabbitmq.NewConn(brokerDsn, rabbitmq.WithConnectionOptionsLogger(log.StandardLogger()))
			if broker != nil && err == nil {
				break
			}

			time.Sleep(2 * time.Second)
		}

		if err != nil || broker == nil {
			return err
		}
	}

	log.Info("RabbitMQ: connection established")

	c.broker = broker

	return nil
}

func (c *Config) closeBroker() error {
	if c.broker != nil {
		if err := c.broker.Close(); err != nil {
			return err
		}

		c.broker = nil
	}

	return nil
}
