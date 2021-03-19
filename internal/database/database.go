package database

import (
	"github.com/golangbb/golangbb/v2/internal"
	"github.com/golangbb/golangbb/v2/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() {
	log.Println("[DATABASE]::CONNECTING 🔌")
	db, err := gorm.Open(sqlite.Open(internal.DATABASENAME), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Println("[DATABASE]::CONNECTION_ERROR 💥")
		log.Fatal(err)
		panic(err)
	}

	DB = db
	log.Println("[DATABASE]::CONNECTED 🔌")
}

func Initialise() {
	log.Println("[DATABASE]::RUNNING_DATABASE_MIGRATIONS 💾")
	err := DB.AutoMigrate(models.Models()...)
	if err != nil {
		log.Println("[DATABASE]::MIGRATION_ERROR 💥")
		log.Fatal(err)
		panic(err)
	}

	log.Println("[DATABASE]::DATABASE_MIGRATIONS_COMPLETE 💾")
}