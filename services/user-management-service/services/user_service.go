package services

import (
	"github.com/google/uuid"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/helpers"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/models"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(phoneNumber, password string) error
	Login(phoneNumber, password string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(phoneNumber, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		ID:          uuid.New().String(),
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
	}
	return s.repo.Create(user)
}

func (s *userService) Login(phoneNumber, password string) (string, error) {
	user, err := s.repo.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	return helpers.GenerateToken(user.ID)
}
