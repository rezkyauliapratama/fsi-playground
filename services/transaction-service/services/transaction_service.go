package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/models"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/repositories"
)

type TransactionService interface {
	CreateDebitTransaction(userID, description string, amount float64, currency string) error
	CreateCreditTransaction(userID, description string, amount float64, currency string) error
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	accountRepo     repositories.AccountRepository
	entryRepo       repositories.EntryRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, accountRepo repositories.AccountRepository, entryRepo repositories.EntryRepository) TransactionService {
	return &transactionService{transactionRepo, accountRepo, entryRepo}
}

func (s *transactionService) Payment(userID, description string, amount float64, currency string) error {

}

func (s *transactionService) CreateDebitTransaction(userID, description string, amount float64, currency string) error {
	transactionID := uuid.New().String()
	debitAccountID := "rezky-account-id" // User's account ID

	// Check if the account has sufficient balance
	balance, err := s.accountRepo.GetBalance(debitAccountID)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.New("insufficient balance")
	}

	transaction := &models.Transaction{
		ID:          transactionID,
		UserID:      userID,
		Amount:      amount,
		Currency:    currency,
		Type:        "debit",
		Description: description,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return err
	}

	debitEntry := &models.Entry{
		ID:            uuid.New().String(),
		TransactionID: transactionID,
		AccountID:     debitAccountID,
		Amount:        amount,
		Type:          "debit",
	}

	creditEntry := &models.Entry{
		ID:            uuid.New().String(),
		TransactionID: transactionID,
		AccountID:     "vendor-account-id", // Vendor's account ID
		Amount:        amount,
		Type:          "credit",
	}

	if err := s.entryRepo.Create(debitEntry); err != nil {
		return err
	}

	if err := s.entryRepo.Create(creditEntry); err != nil {
		return err
	}

	if err := s.accountRepo.UpdateBalance(debitAccountID, -amount); err != nil {
		return err
	}

	if err := s.accountRepo.UpdateBalance("vendor-account-id", amount); err != nil {
		return err
	}

	return nil
}

func (s *transactionService) CreateCreditTransaction(userID, description string, amount float64, currency string) error {
	transactionID := uuid.New().String()
	creditAccountID := "rezky-account-id" // User's account ID

	// Perform the external debit
	externalAccountID := "external-bank-account-id" // The external bank account ID

	transaction := &models.Transaction{
		ID:          transactionID,
		UserID:      userID,
		Amount:      amount,
		Currency:    currency,
		Type:        "credit",
		Description: description,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return err
	}

	creditEntry := &models.Entry{
		ID:            uuid.New().String(),
		TransactionID: transactionID,
		AccountID:     creditAccountID,
		Amount:        amount,
		Type:          "credit",
	}

	debitEntry := &models.Entry{
		ID:            uuid.New().String(),
		TransactionID: transactionID,
		AccountID:     externalAccountID,
		Amount:        amount,
		Type:          "debit",
	}

	if err := s.entryRepo.Create(creditEntry); err != nil {
		return err
	}

	if err := s.entryRepo.Create(debitEntry); err != nil {
		return err
	}

	if err := s.accountRepo.UpdateBalance(creditAccountID, amount); err != nil {
		return err
	}

	return nil
}
