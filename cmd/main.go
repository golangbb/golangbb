package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	constants "github.com/golangbb/golangbb/v2/internal"
	"github.com/golangbb/golangbb/v2/internal/database"
	"github.com/golangbb/golangbb/v2/internal/models"
)


func init() {
	fmt.Println("[INIT]::INITIALISING 🏗️")

	database.Connect()

	fmt.Println("[INIT]::RUNNING_DATABASE_MIGRATIONS 💾")
	database.DB.AutoMigrate(&models.User{})

	fmt.Println("[INIT]::INITIALISATION_COMPLETE 🏗️")
}

func main() {
	fmt.Println("[MAIN]::BOOTSTRAPPING 🚀")
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		err := c.SendStatus(200)
		if err != nil {
			return err
		}

		return c.SendString("ok")
	})

	app.Listen(":" + constants.PORT)
	fmt.Println("[MAIN]::BOOTSTRAPPED 🚀")
}
