package models

import (
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
)

type CommonFields struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type UserSerializer struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	Currency    string  `json:"currency"`
	TotalIncome float64 `json:"total_income"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type IncomeSerializer struct {
	CommonFields
	Amount       float64               `json:"amount"`
	Category     *string               `json:"category"`
	DateRecieved time.Time             `json:"date_recieved"`
	Type         utils.IncomeType      `json:"type"`
	Frequency    utils.IncomeFrequency `json:"frequency"`
}

type AccountSerializer struct {
	CommonFields
	Name     string  `json:"name"`
	Currency string  `json:"currency"`
	IsActive bool    `json:"is_active"`
	Balance  float64 `json:"balance"`
}

func SerializeUser(user *User) *UserSerializer {
	return &UserSerializer{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		Currency:    user.Currency,
		TotalIncome: user.TotalIncome,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func SerializeIncome(income *Income) *IncomeSerializer {
	return &IncomeSerializer{
		CommonFields: CommonFields{
			ID:        income.ID.String(),
			CreatedAt: income.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: income.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Amount:       income.Amount,
		Category:     income.Category,
		DateRecieved: income.DateRecieved,
		Type:         income.Type,
		Frequency:    income.Frequency,
	}
}

func SerializeAccount(account *Account) *AccountSerializer {
	return &AccountSerializer{
		CommonFields: CommonFields{
			ID:        account.ID.String(),
			CreatedAt: account.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: account.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Name:     account.Name,
		Currency: account.Currency,
		IsActive: account.IsActive,
		Balance:  account.Balance,
	}
}
