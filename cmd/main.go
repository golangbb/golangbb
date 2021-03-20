package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golangbb/golangbb/v2/internal"
	"github.com/golangbb/golangbb/v2/internal/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)


func initialise() *sql.DB {
	log.Println("[INIT]::INITIALISING 🏗️")
	dbConnection := database.Connect(sqlite.Open(internal.DATABASENAME), gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	sqlDb, dbErr := dbConnection.DB()
	if dbErr != nil {
		log.Println("[INIT]::CONNECTION_ERROR 💥")
		log.Fatal(dbErr)
		panic(dbErr)
	}
	database.Initialise()
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
