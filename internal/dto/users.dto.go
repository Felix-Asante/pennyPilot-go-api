package dto

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email" errormgs:"Invalid email address"`
	Password string `json:"password" validate:"required,min=8" errormgs:"Password must be at least 8 characters long"`
	FullName string `json:"full_name" validate:"required" errormgs:"Full name is required"`
	Currency string `json:"currency" validate:"required" errormgs:"Currency is required"`
}

type LoginDto struct {
	Email    string `json:"email" validate:"required,email" errormgs:"Invalid email address"`
	Password string `json:"password" validate:"required,min=8" errormgs:"Password must be at least 8 characters long"`
}

type ForgotPasswordDto struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDto struct {
	ResetToken  string `json:"reset_token" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
