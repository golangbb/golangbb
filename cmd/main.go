package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	constants "github.com/golangbb/golangbb/v2/internal"
	"github.com/golangbb/golangbb/v2/internal/database"
)


func init() {
	fmt.Println("[INIT]::INITIALISING 🏗️")
	database.Connect()
	database.Initialise()
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
