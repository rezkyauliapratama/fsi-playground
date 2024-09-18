package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/services"
)

func RegisterTransactionHandlers(app *fiber.App, service services.TransactionService) {
	app.Post("/transactions/debit", func(c *fiber.Ctx) error {
		var request struct {
			UserID        string  `json:"user_id"`
			Description   string  `json:"description"`
			CreditAccount string  `json:"credit_account"`
			Amount        float64 `json:"amount"`
			Currency      string  `json:"currency"`
			Timestamp     string  `json:"timestamp"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := service.CreateDebitTransaction(request.UserID, request.CreditAccount, request.Description, request.Amount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Debit transaction created"})
	})

	app.Post("/transactions/credit", func(c *fiber.Ctx) error {
		var request struct {
			UserID      string  `json:"user_id"`
			Description string  `json:"description"`
			Amount      float64 `json:"amount"`
			Currency    string  `json:"currency"`
			Timestamp   string  `json:"timestamp"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := service.CreateCreditTransaction(request.UserID, request.Description, request.Amount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Credit transaction created"})
	})
}
