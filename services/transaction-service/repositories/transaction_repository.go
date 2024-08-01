package repositories

import (
	"database/sql"
	"time"

	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	Update(transaction *models.Transaction) error
	GetByID(transactionID string) (*models.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	query := "INSERT INTO transactions (id, user_id, amount, currency, type, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())"
	_, err := r.db.Exec(query, transaction.ID, transaction.UserID, transaction.Amount, transaction.Currency, transaction.Type, transaction.Description)
	return err
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	query := "UPDATE transactions SET amount = ?, currency = ?, type = ?, description = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.Exec(query, transaction.Amount, transaction.Currency, transaction.Type, transaction.Description, transaction.ID)
	return err
}

func (r *transactionRepository) GetByID(transactionID string) (*models.Transaction, error) {
	var transaction models.Transaction
	var createdAt, updatedAt []uint8
	query := "SELECT id, user_id, amount, currency, type, description, created_at, updated_at FROM transactions WHERE id = ?"
	err := r.db.QueryRow(query, transactionID).Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Currency, &transaction.Type, &transaction.Description, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	transaction.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, err
	}
	transaction.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
