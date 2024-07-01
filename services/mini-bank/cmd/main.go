package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rezkyauliapratama/fsi-playground/services/mini-bank/handler"
)

func setupRoutes(app *fiber.App) {
	// Initialize default config
	app.Post("/api/v1/transaction", handler.PostTransaction)
}
func main() {
	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
