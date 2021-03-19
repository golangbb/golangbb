package database

import (
	"fmt"
	constants "github.com/golangbb/golangbb/v2/internal"
	"github.com/golangbb/golangbb/v2/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	fmt.Println("[DATABASE]::CONNECTING 🔌")
	db, err := gorm.Open(sqlite.Open(constants.DATABASE_NAME), &gorm.Config{})
	if err != nil {
		fmt.Println("[DATABASE]::CONNECTION_ERROR 💥")
		panic(err)
	}

	DB = db
	fmt.Println("[DATABASE]::CONNECTED 🔌")
}

func Initialise() {
	fmt.Println("[DATABASE]::RUNNING_DATABASE_MIGRATIONS 💾")
	DB.AutoMigrate(&models.User{})
	fmt.Println("[DATABASE]::DATABASE_MIGRATIONS_COMPLETE 💾")
}