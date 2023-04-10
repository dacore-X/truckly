package dto

// UserSignUpRequestBody represents the request body with data
// sent by the user to API to sign up in the application
type UserSignUpRequestBody struct {
	Surname     string `json:"surname" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
	IsCourier   bool   `json:"is_courier"`
}

// UserLoginRequestBody represents the request body with data
// sent by the user to API to log in to the application
type UserLoginRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserBanParams represents URI with user's ID
// to identify user to be banned/unbanned
type UserBanParams struct {
	ID int `uri:"id" binding:"required,min=1"`
}
