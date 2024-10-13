package main

import (
	"internal-iam-service/handlers"
	"internal-iam-service/repositories"
	"internal-iam-service/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func main() {
	// Initialize a new Fiber instance
	app := fiber.New(fiber.Config{
		Prefork:       true, // Optional: Enable Prefork if you want to scale with multiple processes
		CaseSensitive: true, // Routes are case-sensitive
		StrictRouting: true, // Strict routing ensures /foo and /foo/ are treated differently
	})

	// Add middleware
	app.Use(recover.New()) // Automatically recover from panics

	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize repository, service, and handler
	ketoRepo := repositories.NewKetoRepository(logger)
	ketoService := services.NewKetoService(ketoRepo, logger)
	ketoHandler := handlers.NewKetoHandler(ketoService, logger)

	// Define routes
	app.Get("/users/:userID/modules", ketoHandler.GetUserModules)

	// Start the Fiber app
	if err := app.Listen(":8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
