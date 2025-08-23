package dto

type CreateAccountDto struct {
	Name     string `json:"name" validate:"required,min=3,max=20"`
	Currency string `json:"currency" validate:"required"`
}

type UpdateAccountDto struct {
	Name     *string `json:"name" validate:"omitempty"`
	Currency *string `json:"currency" validate:"omitempty"`
	IsActive *bool   `json:"is_active" validate:"omitempty"`
}
