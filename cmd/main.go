package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		err := c.SendStatus(200)
		if err != nil {
			return err
		}

		return c.SendString("ok")
	})

	app.Listen(":3000")
}