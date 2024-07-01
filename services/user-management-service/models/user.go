package models

import "time"

type User struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
