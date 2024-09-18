package repositories

import (
	"database/sql"
	"time"

	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByPhoneNumber(phoneNumber string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {

	query := "INSERT INTO users (id, phone_number, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
	_, err := r.db.Exec(query, user.ID, user.PhoneNumber, user.Email, user.Password)
	return err
}

func (r *userRepository) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	var createdAt, updatedAt []uint8
	query := "SELECT id, phone_number, password, created_at, updated_at FROM users WHERE phone_number = ?"
	err := r.db.QueryRow(query, phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, err
	}
	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
