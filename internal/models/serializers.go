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

type AllocationRuleSerializer struct {
	CommonFields
	Strategy utils.AllocationStrategy `json:"strategy"`
	Value    float64                  `json:"value"`
	TargetID string                   `json:"target_id"`
	Active   bool                     `json:"active"`
}

type EnvelopeSerializer struct {
	CommonFields
	Name           string                    `json:"name"`
	CurrentAmount  float64                   `json:"current_amount"`
	TargetAmount   float64                   `json:"target_amount"`
	AutoAllocate   bool                      `json:"auto_allocate"`
	TargetedDate   string                    `json:"targeted_date"`
	IsActive       bool                      `json:"is_active"`
	Account        *AccountSerializer        `json:"account"`
	AllocationRule *AllocationRuleSerializer `json:"allocation_rule"`
}

func SerializeEnvelope(envelope *Envelope) *EnvelopeSerializer {
	if envelope == nil {
		return nil
	}

	serializedEnvelope := EnvelopeSerializer{
		CommonFields: CommonFields{
			ID:        envelope.ID.String(),
			CreatedAt: envelope.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: envelope.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Name:           envelope.Name,
		CurrentAmount:  envelope.CurrentAmount,
		TargetAmount:   envelope.TargetAmount,
		AutoAllocate:   envelope.AutoAllocate,
		IsActive:       envelope.IsActive,
		Account:        SerializeAccount(&envelope.Account),
		AllocationRule: SerializeAllocationRule(envelope.AllocationRule),
	}

	if envelope.TargetedDate != nil {
		serializedEnvelope.TargetedDate = envelope.TargetedDate.Format("2006-01-02 15:04:05")
	}

	return &serializedEnvelope
}

func SerializeUser(user *User) *UserSerializer {
	if user == nil {
		return nil
	}
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
	if income == nil {
		return nil
	}
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
	if account == nil {
		return nil
	}
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

func SerializeAllocationRule(allocationRule *AllocationRule) *AllocationRuleSerializer {
	if allocationRule == nil {
		return nil
	}
	return &AllocationRuleSerializer{
		CommonFields: CommonFields{
			ID:        allocationRule.ID.String(),
			CreatedAt: allocationRule.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: allocationRule.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Strategy: allocationRule.Strategy,
		Value:    allocationRule.Value,
		TargetID: allocationRule.TargetID.String(),
		Active:   allocationRule.Active,
	}
}
