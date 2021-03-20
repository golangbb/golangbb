package database

import (
	"github.com/golangbb/golangbb/v2/internal/models"
	"gorm.io/gorm"
	"log"
)

var DBConnection *gorm.DB

func Connect(dialector gorm.Dialector, config gorm.Config) *gorm.DB {
	log.Println("[DATABASE]::CONNECTING 🔌")
	db, err := gorm.Open(dialector, &config)
	if err != nil {
		log.Println("[DATABASE]::CONNECTION_ERROR 💥")
		log.Fatal(err)
		panic(err)
	}

	DBConnection = db
	log.Println("[DATABASE]::CONNECTED 🔌")
	return DBConnection
}

func Initialise() {
	log.Println("[DATABASE]::RUNNING_DATABASE_MIGRATIONS 💾")
	err := DBConnection.AutoMigrate(models.Models()...)
	if err != nil {
		log.Println("[DATABASE]::MIGRATION_ERROR 💥")
		log.Fatal(err)
		panic(err)
	}

	log.Println("[DATABASE]::DATABASE_MIGRATIONS_COMPLETE 💾")
}