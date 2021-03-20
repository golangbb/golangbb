package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golangbb/golangbb/v2/internal"
	"github.com/golangbb/golangbb/v2/internal/database"
	"github.com/golangbb/golangbb/v2/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)


func initialise() *sql.DB {
	log.Println("[INIT]::INITIALISING 🏗️")
	dbConnection, err := database.Connect(sqlite.Open(internal.DATABASENAME), gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Println("[INIT]::CREATE_CONNECTION_ERROR 💥")
		log.Fatal(err)
		panic(err)
	}

	sqlDb, err := dbConnection.DB()
	if err != nil {
		log.Println("[INIT]::GET_UNDERLYING_SQL_CONNECTION_ERROR 💥")
		log.Fatal(err)
		panic(err)
	}

	err = database.Initialise(models.Models()...)
	if err != nil {
		log.Println("[INIT]::DATABASE_INITIALISE_ERROR 💥")
		log.Fatal(err)
		panic(err)
	}

	log.Println("[INIT]::INITIALISATION_COMPLETE 🏗️")
	return sqlDb
}

func main() {
	db := initialise()
	defer db.Close()

	log.Println("[MAIN]::BOOTSTRAPPING 🚀")
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		err := c.SendStatus(200)
		if err != nil {
			return err
		}

		return c.SendString("ok")
	})

	log.Println("[MAIN]::BOOTSTRAPPED 🚀")
	log.Fatal(app.Listen(":" + internal.PORT))
}
