package models

type RegisterReq struct {
	Username       string `json:"username" validate:"required,max=40" errormgs:"invalid username: can not be empty or more than 40 characters"`
	Password       string `json:"password" validate:"required,max=40,min=7" errormgs:"invalid password: can not be less than 7 or greater than 40 characters"`
	Email          string `json:"email" validate:"required,min=5,email" errormgs:"invalid email: can not be less than 5 characters"`
	PassportNumber string `json:"passportNumber" validate:"required,min=9" errormgs:"invalid passport number: can not be less than 9 characters"`
	PaymentSource  string `json:"paymentSource" validate:"required,min=12" errormgs:"invalid card number: can not be less than 12 characters"`
	PhoneNumber    string `json:"phoneNumber" validate:"required,min=7" errormgs:"invalid phone number: can not be less than 7 characters"`
}

type RegisterRes struct {
	UUID string `json:"uuid"`
}

type LoginReq struct {
	Password string `json:"password" validate:"required,max=40,min=7" errormgs:"invalid password: can not be less than 7 or greater than 40 characters"`
	Email    string `json:"email" validate:"required,min=5" errormgs:"invalid email: can not be less than 5 characters"`
}

type ResetPasswordReq struct {
	OldPassword string `json:"oldPassword" validate:"required,min=7,max=40"`
	NewPassword string `json:"newPassword" validate:"required,min=7,max=40"`
	Token       string
}
