package dto

import "time"

// RoleMeta represents struct
// with meta data about user's role
type RoleMeta struct {
	IsAdmin   bool `json:"is_admin"`
	IsCourier bool `json:"is_courier"`
	IsBanned  bool `json:"is_banned"`
}

// UserMeResponse represents the response body
// with user's account data sent to the user by API
type UserMeResponse struct {
	ID          int       `json:"id"`
	Surname     string    `json:"surname"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	Meta        *RoleMeta `json:"meta"`
}

// UserInfoResponse represents the response body
// with private user's data sent to the user
// verification handlers by API
type UserInfoResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserMetaResponse represents the response body
// with metadata about user
type UserMetaResponse struct {
	UserID    int     `json:"user_id"`
	IsAdmin   bool    `json:"is_admin"`
	IsCourier bool    `json:"is_courier"`
	IsBanned  bool    `json:"is_banned"`
	Rating    float32 `json:"rating"`
}
