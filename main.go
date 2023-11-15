package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {

	//? Connect to database
	db := ConnectToDB()

	//? App initialization
	app := fiber.New()

	//? Routes
	UserRoutes(app, db)

	//? 404 Error handling
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})

	//? Server setup
	app.Listen(":3000")
}
