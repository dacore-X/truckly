package dto

import "time"

type UserResponseMeBody struct {
	ID          int64     `json:"id"`
	Surname     string    `json:"surname"`
	Name        string    `json:"name"`
	Patronymic  string    `json:"patronymic"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserResponseInfoBody struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
