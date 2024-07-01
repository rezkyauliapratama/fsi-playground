package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/helpers"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/services"
)

func RegisterUserHandlers(app *fiber.App, service services.UserService) {
	app.Post("/register", func(c *fiber.Ctx) error {
		var request struct {
			PhoneNumber string `json:"phone_number"`
			Password    string `json:"password"`
		}
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := service.Register(request.PhoneNumber, request.Password); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created"})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var request struct {
			PhoneNumber string `json:"phone_number"`
			Password    string `json:"password"`
		}
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		token, err := service.Login(request.PhoneNumber, request.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"token": token})
	})

	app.Get("/protected", AuthenticateJWT, func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		return c.JSON(fiber.Map{"message": "You have accessed a protected route", "userID": userID})
	})
}

func AuthenticateJWT(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	token, err := helpers.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("userID", claims["userID"])
		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
}
