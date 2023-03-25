package dto

import "time"

// UserMeResponse represents the response body
// with user's account data sent to the user by API
type UserMeResponse struct {
	ID          int64     `json:"id"`
	Surname     string    `json:"surname"`
	Name        string    `json:"name"`
	Patronymic  string    `json:"patronymic"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserInfoResponse represents the response body
// with private user's data sent to the user
// verification handlers by API
type UserInfoResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
