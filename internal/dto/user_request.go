package dto

type UserRequestSignUpBody struct {
	Surname     string `json:"surname"`
	Name        string `json:"name"`
	Patronymic  string `json:"patronymic"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	IsCourier   bool   `json:"is_courier"`
}

type UserRequestLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
