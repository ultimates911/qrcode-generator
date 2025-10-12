package dto

type RegisterRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	SecondPassword  string `json:"second_password" validate:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}