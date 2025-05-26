package entity

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDb(gormDB *gorm.DB) {
	db = gormDB
}

func Db() *gorm.DB {
	if db == nil {
		log.Panic("entity: db connection is not initialized")
	}

	return db
}
