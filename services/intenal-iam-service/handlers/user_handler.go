package handlers

import (
	"internal-iam-service/services"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler defines the handler for user operations with logger
type UserHandler struct {
	userService *services.KetoService
	logger      *zap.Logger
}

// NewUserHandler initializes a new UserHandler with logger
func NewUserHandler(userService *services.KetoService, logger *zap.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

// GetUserModulesHandler handles the request to get the modules a user can access
func (h *UserHandler) GetUserModulesHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		h.logger.Error("Missing userID in request", zap.String("request_id", c.Locals("requestid").(string)))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	h.logger.Info("Fetching modules for user", zap.String("userID", userID))

	// Call the service to get the list of modules the user has access to
	modulesByAction, err := h.userService.GetUserModules(c.Context(), userID)
	if err != nil {
		h.logger.Error("Error fetching user modules", zap.String("userID", userID), zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.logger.Info("Successfully fetched modules for user", zap.String("userID", userID), zap.Any("modules", modulesByAction))

	// Return the modules as a JSON response
	return c.JSON(fiber.Map{
		"userID":          userID,
		"modulesByAction": modulesByAction,
	})
}
