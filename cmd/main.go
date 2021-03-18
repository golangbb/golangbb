package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golangbb/golangbb/v2/pkg/helpers"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

var (
	PORT = helpers.GetEnv("PORT", "3000")
	DATABASE_NAME = helpers.GetEnv("DATABASE_NAME", "golangbb.db")
)

func init() {
	fmt.Println("🏗️ initializing golangBB...")

	fmt.Println("💾 connecting to database...")
	db, err := gorm.Open(sqlite.Open(DATABASE_NAME), &gorm.Config{})
	if err != nil {
		panic("💥 failed to connect database")
	}

	fmt.Println("💾 running database migrations...")
	db.AutoMigrate()
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
