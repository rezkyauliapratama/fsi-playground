package main

import (
	"internal-iam-service/handlers"
	"internal-iam-service/repositories"
	"internal-iam-service/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Initialize Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync() // Flushes the logger buffer, should be called at the end

	// Initialize the Fiber web application
	app := fiber.New()

	// Initialize repository, service, and handler layers with logger
	ketoRepo := repositories.NewKetoRepository(logger)
	ketoService := services.NewKetoService(ketoRepo, logger)
	userHandler := handlers.NewUserHandler(ketoService, logger)

	// Define routes for the API
	app.Get("/users/:id/modules", userHandler.GetUserModulesHandler)

	// Start the Fiber app on the specified port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	logger.Info("Starting server", zap.String("port", port))
	log.Fatal(app.Listen(":" + port))
}
