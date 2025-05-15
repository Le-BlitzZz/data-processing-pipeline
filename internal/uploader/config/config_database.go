package config

import (
	"Le-BlitzZz/streaming-etl-app/internal/database"

	log "github.com/sirupsen/logrus"
)

const (
	dbServer   = "mariadb:4001"
	dbName     = "etlstream"
	dbUser     = "root"
	dbPassword = "etlstream"
	dbTimeout  = 15
)

func (c *Config) Db() *database.Database {
	if c.db == nil {
		log.Fatal("config: database not connected")
	}

	return c.db
}
