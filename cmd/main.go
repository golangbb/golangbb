package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golangbb/golangbb/v2/pkg/helpers"
)

var (
	PORT = helpers.GetEnv("PORT", "3000")
)

func init() {
	fmt.Println("🏗️ initializing golangBB...")
}

func main() {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		err := c.SendStatus(200)
		if err != nil {
			return err
		}

		return c.SendString("ok")
	})

	app.Listen(":" + PORT)
}
