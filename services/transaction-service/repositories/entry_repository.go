package repositories

import (
	"database/sql"
	"time"

	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/models"
)

type EntryRepository interface {
	Create(entry *models.Entry) error
	GetByTransactionID(transactionID string) ([]*models.Entry, error)
}

type entryRepository struct {
	db *sql.DB
}

func NewEntryRepository(db *sql.DB) EntryRepository {
	return &entryRepository{db}
}

func (r *entryRepository) Create(entry *models.Entry) error {
	query := "INSERT INTO entries (id, transaction_id, account_id, amount, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())"
	_, err := r.db.Exec(query, entry.ID, entry.TransactionID, entry.AccountID, entry.Amount, entry.Type)
	return err
}

func (r *entryRepository) GetByTransactionID(transactionID string) ([]*models.Entry, error) {
	var entries []*models.Entry
	query := "SELECT id, transaction_id, account_id, amount, type, created_at, updated_at FROM entries WHERE transaction_id = ?"
	rows, err := r.db.Query(query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.Entry
		var createdAt, updatedAt []uint8
		if err := rows.Scan(&entry.ID, &entry.TransactionID, &entry.AccountID, &entry.Amount, &entry.Type, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		entry.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, err
		}
		entry.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}
