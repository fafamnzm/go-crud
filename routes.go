package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserRoutes(router fiber.Router, db *gorm.DB) {
	//? Check server connection
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("You are connected")
	})

	//? Get all users
	router.Get("/users", func(c *fiber.Ctx) error {
		var users []User
		result := db.Find(&users)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(users)
	})

	//? Create a user
	router.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return err
		}
		result := db.Create(&user)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(user)
	})

	//? Get a single user
	router.Get("/users/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}
		var user User
		db.First(&user, id)
		return c.JSON(user)
	})

	//? Update a specific user
	router.Put("/users/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return err
		}
		result := db.Model(&User{}).Where("id = ?", id).Updates(user)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(user)
	})

	//? Delete a specific user
	router.Delete("/users/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}
		db.Delete(&User{}, id)
		return c.SendString("User deleted successfully")
	})
}
