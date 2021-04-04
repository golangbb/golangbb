package database

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

var DBConnection *gorm.DB

var NoDatabaseConnectionErr = errors.New("no database connection")

func Connect(dialector gorm.Dialector, config gorm.Config) (*gorm.DB, error) {
	log.Println("[DATABASE]::CONNECTING 🔌")
	db, err := gorm.Open(dialector, &config)
	if err != nil {
		log.Println("[DATABASE]::CONNECTION_ERROR 💥")
		return nil, err
	}

	DBConnection = db
	log.Println("[DATABASE]::CONNECTED 🔌")
	return DBConnection, nil
}

func Initialise(models ...interface{}) error {
	log.Println("[DATABASE]::RUNNING_DATABASE_MIGRATIONS 💾")
	if DBConnection == nil {
		return NoDatabaseConnectionErr
	}

	err := DBConnection.AutoMigrate(models...)
	if err != nil {
		log.Println("[DATABASE]::MIGRATION_ERROR 💥")
		return err
	}

	log.Println("[DATABASE]::DATABASE_MIGRATIONS_COMPLETE 💾")
	return nil
}
