package entity

import "time"

// User represents the user's entity data struct
type User struct {
	ID          int       `json:"id"`
	Surname     string    `json:"surname"`
	Name        string    `json:"name"`
	Patronymic  string    `json:"patronymic"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
}
