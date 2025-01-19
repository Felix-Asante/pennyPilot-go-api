package repositories

import "time"

type CreateIncomeDto struct {
	Amount       float64    `json:"amount" validate:"required;min=1"`
	DateReceived *time.Time `json:"date_received" validate:"required"`
	Type         string     `json:"type" validate:"required"`
	Frequency    string     `json:"frequency" validate:"required"`
	User         string
}

type UpdateIncomeDto struct {
	Amount       *float64         `json:"amount" validate:"omitempty,gt=0"`
	DateReceived *time.Time       `json:"date_received" validate:"omitempty,datetime"`
	Type         *IncomeType      `json:"type" validate:"omitempty,oneof=salary investment freelance others"`
	Frequency    *IncomeFrequency `json:"frequency" validate:"omitempty,oneof=weekly monthly yearly bi-weekly one-time"`
}
