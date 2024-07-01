package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	error_enum "github.com/rezkyauliapratama/fsi-playground/services/mini-bank/helper"
)

type Transaction struct {
	// dbConn *sql.DB
}

type TransactionRequest struct {
	TransactionId string          `json:"transaction_id"`
	AccountId     string          `json:"account_id"`
	Type          TransactionType `json:"transaction_type"`
	Amount        int64           `json:"amount"`
	Timestamp     uint64          `json:"timestamp"`
	Description   string          `json:"description"`
}

type TransactionType string

const (
	Withdrawal TransactionType = "WITHDRAWAL"
	Transfer   TransactionType = "TRANSFER"
	Saving     TransactionType = "SAVING"
)

//	func NewTransactionHandler(dbConn *sql.DB) Transaction {
//		return Transaction{dbConn: dbConn}
//	}
func NewTransactionHandler() Transaction {
	return Transaction{}
}

func (t TransactionType) Validate() error {
	switch t {
	case Withdrawal, Transfer, Saving:
		return nil
	}
	return error_enum.ErrTransactionTypeNotFound
}

func PostTransaction(c *fiber.Ctx) error {
	payload := TransactionRequest{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(error_enum.ErrTransactionTypeNotFound)
	}

	return c.Status(http.StatusAccepted).JSON(payload)

}
