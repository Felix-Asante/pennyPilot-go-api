package repositories

import "time"

type CreateIncomeDto struct {
	Amount       float64    `json:"amount" validate:"required;min=1"`
	DateReceived *time.Time `json:"date_received" validate:"required"`
	Type         string     `json:"type" validate:"required"`
	Frequency    string     `json:"frequency" validate:"required"`
	User         string
}
