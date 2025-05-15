package config

import (
	"Le-BlitzZz/streaming-etl-app/internal/broker"
	"Le-BlitzZz/streaming-etl-app/internal/database"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	mb *broker.MessageBroker
	db *database.Database

	rawExchange       string
	processedExchange string
	rawQueue          string
	processedQueue    string
}

func NewConfig() (*Config, error) {
	c := &Config{
		rawExchange:       rawExchange,
		processedExchange: processedExchange,
		rawQueue:          rawQueue,
		processedQueue:    processedQueue,
	}

	mbConfig := broker.NewConfig(mbUser, mbPassword, mbServer)

	mb, err := broker.NewMessageBroker(mbConfig)
	if err != nil {
		return nil, err
	}

	c.mb = mb

	dbConfig := database.NewConfig(dbUser, dbPassword, dbServer, dbName, dbTimeout)

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		c.shutdownMessageBroker()
		return nil, err
	}

	c.db = db

	log.Debug("config: successfully initialized")

	return c, nil
}

func (c *Config) Shutdown() {
	c.shutdownMessageBroker()

	if err := c.db.Close(); err != nil {
		log.Errorf("could not close database connection: %s", err)
	} else {
		log.Debug("closed database connection")
	}
}

func (c *Config) shutdownMessageBroker() {
	if err := c.mb.Close(); err != nil {
		log.Errorf("could not close message broker connection: %s", err)
	} else {
		log.Debug("closed message broker connection")
	}
}
