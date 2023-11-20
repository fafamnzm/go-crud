package main

import (
	"fmt"
	"strconv"
	"strings"

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

	//? Generate OTP and email or sms it
	router.Post("/generateOTP", func(c *fiber.Ctx) error {
		//? Parse request body into User struct
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return err
		}

		//? Generate a random number
		otpNumber := strconv.Itoa(generateRandomNumber())

		//? Create a new instance of the OTPVerification struct
		otpVerification := &OTPVerification{
			Email: user.Email,
			OTP:   otpNumber,
		}

		if err := setOtpVerification(CommonCache, user.Email, otpVerification); err != nil {
			return err
		}

		//? Print the token to the console similar to emailing or sending an sms
		fmt.Println("Generated Token:", otpVerification)

		//? Return the random number in the response
		return c.JSON(fiber.Map{
			"message": "OTP number sent",
		})
	})

	router.Post("/verifyOTP", func(c *fiber.Ctx) error {
		//? Parse request body into VerificationRequest struct
		request := new(OTPVerification)
		if err := c.BodyParser(request); err != nil {
			return err
		}

		otpVerification, err := getOtpVerification(CommonCache, request.Email)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
		}

		//? Compare the received OTP with the stored OTP
		if request.Email == otpVerification.Email && request.OTP == otpVerification.OTP {
			//? Check newUser in db and create it if not
			var newUser User
			if err := db.Where("email = ?", request.Email).First(&newUser).Error; err != nil {
				//? In case the user doesn't exist, add it to the database
				newUser = User{
					Email: request.Email,
					//? Set other user fields as needed
				}
				if err := db.Create(&newUser).Error; err != nil {
					//? Handle the error
					return err
				}
			}

			//? Generate JWT token
			tokenString, err := generateJWTToken(newUser)
			if err != nil {
				return err
			}

			//? Send the JWT token in the response
			return c.JSON(fiber.Map{
				"token": tokenString,
			})
		} else {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Invalid OTP"})
		}
	})

	router.Get("/me", func(c *fiber.Ctx) error {
		//? Extract bearer token from request headers
		authHeader := c.Get("Authorization")
		tokenString := strings.Split(authHeader, "Bearer ")[1]

		//? Verify the JWT token
		claims, err := verifyJWTToken(tokenString)
		if err != nil {
			//? Invalid token handler
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		//? Get the user's id from the claims
		id := claims.ID

		//? Retrieve user from the database using the id
		var user User
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
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
