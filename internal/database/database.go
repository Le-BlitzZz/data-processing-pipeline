package database

import (
	"Le-BlitzZz/streaming-etl-app/internal/entity"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	connection *gorm.DB
	config     *Config
}

func NewDatabase(config *Config) (*Database, error) {
	db := &Database{
		config: config,
	}

	if err := db.connect(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) Close() error {
	if db.connection != nil {
		sqlDb, err := db.connection.DB()
		if err != nil {
			return err
		}

		return sqlDb.Close()
	}

	return nil
}

func (db *Database) Init() {
	for name, entity := range entity.Entities {
		if err := db.connection.AutoMigrate(entity); err != nil {
			log.Errorf("Failed migrating %s", name)
		}
	}
}

func (db *Database) connect() error {
	dbDsn := db.config.DatabaseDsn()
	connection, err := gorm.Open(mysql.Open(dbDsn), &gorm.Config{})

	if err != nil || connection == nil {
		log.Println("MariaDB: waiting to become available")

		for i := 1; i <= 12; i++ {
			connection, err = gorm.Open(mysql.Open(dbDsn), &gorm.Config{})

			if connection != nil && err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if err != nil || connection == nil {
			return err
		}
	}

	log.Println("MariaDB: connection established")

	db.connection = connection

	return nil
}

func (db *Database) Create(value any) error {
	return db.connection.Create(value).Error
}
