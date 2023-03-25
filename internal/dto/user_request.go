package dto

// UserSignUpRequestBody represents the request body with data
// sent by the user to API to sign up in the application
type UserSignUpRequestBody struct {
	Surname     string `json:"surname"`
	Name        string `json:"name"`
	Patronymic  string `json:"patronymic"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	IsCourier   bool   `json:"is_courier"`
}

// UserLoginRequestBody represents the request body with data
// sent by the user to API to log in to the application
type UserLoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
