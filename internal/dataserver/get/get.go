package get

import (
	"Le-BlitzZz/streaming-etl-app/internal/dataserver/config"

	log "github.com/sirupsen/logrus"
)

var conf *config.Config

func SetConfig(c *config.Config) {
	if c == nil {
		log.Panic("get.SetConfig: config cannot be nil")
	}

	conf = c
}

func Config() *config.Config {
	if conf == nil {
		log.Panic("get.Config: config is not set")
	}

	return conf
}