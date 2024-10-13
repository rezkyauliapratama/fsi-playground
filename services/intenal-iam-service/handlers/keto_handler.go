package handlers

import (
	"context"
	"internal-iam-service/services"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// KetoHandler handles the HTTP request related to Keto service
type KetoHandler struct {
	service *services.KetoService
	logger  *zap.Logger
}

// NewKetoHandler initializes a new KetoHandler
func NewKetoHandler(service *services.KetoService, logger *zap.Logger) *KetoHandler {
	return &KetoHandler{service: service, logger: logger}
}

// GetUserModules handles the request to fetch accessible modules for a user
func (h *KetoHandler) GetUserModules(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		h.logger.Error("Missing userID in request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing userID",
		})
	}

	// Set a timeout for the request context (5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch accessible modules for the user
	h.logger.Info("Fetching accessible modules", zap.String("userID", userID))
	modules, err := h.service.GetUserModules(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to fetch modules for user", zap.String("userID", userID), zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch modules",
		})
	}

	// Return the modules and actions
	h.logger.Info("Successfully fetched modules", zap.String("userID", userID))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"userID":  userID,
		"modules": modules,
	})
}
