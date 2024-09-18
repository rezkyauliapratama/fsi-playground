package repositories

import (
	"database/sql"
)

type AccountRepository interface {
	GetBalance(accountID string) (float64, error)
	UpdateBalance(accountID string, amount float64) error
	GetAccountByUser(userID string, accountType string) (string, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) GetBalance(accountID string) (float64, error) {
	var balance float64
	query := "SELECT balance FROM accounts WHERE id = ?"
	err := r.db.QueryRow(query, accountID).Scan(&balance)
	return balance, err
}

func (r *accountRepository) UpdateBalance(accountID string, amount float64) error {
	query := "UPDATE accounts SET balance = balance + ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.Exec(query, amount, accountID)
	return err
}

func (r *accountRepository) GetAccountByUser(userID string, accountType string) (string, error) {
	var accountId string
	query := "Select id from accounts WHERE user_id = ? and type = ?"
	err := r.db.QueryRow(query, userID, accountType).Scan(&accountId)
	return accountId, err
}
