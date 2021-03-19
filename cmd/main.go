package main

import (
	"github.com/gofiber/fiber/v2"
	constants "github.com/golangbb/golangbb/v2/internal"
	"github.com/golangbb/golangbb/v2/internal/database"
	"log"
)


func init() {
	log.Println("[INIT]::INITIALISING 🏗️")
	database.Connect()
	database.Initialise()
	log.Println("[INIT]::INITIALISATION_COMPLETE 🏗️")
}

func main() {
	log.Println("[MAIN]::BOOTSTRAPPING 🚀")
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		err := c.SendStatus(200)
		if err != nil {
			return err
		}

		return c.SendString("ok")
	})

	log.Fatal(app.Listen(":" + constants.PORT))
}
