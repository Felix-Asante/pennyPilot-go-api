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

type CreateGoalDto struct {
	Name            string     `json:"name" validate:"required,min=2"`
	TargetAmount    float64    `json:"target_amount" validate:"required,min=0.00"`
	AllocationPoint float64    `json:"allocation_point" validate:"required,min=1,max=100"`
	DueDate         *time.Time `json:"due_date" validate:"required"`
	Account         string     `json:"account_id" validate:"required"`
}
type UpdateGoalDto struct {
	Name            string     `json:"name" validate:"omitempty,min=2"`
	TargetAmount    float64    `json:"target_amount" validate:"omitempty,gt=0"`
	AllocationPoint float64    `json:"allocation_point" validate:"omitempty,min=1,max=100"`
	DueDate         *time.Time `json:"due_date" validate:"omitempty"`
	Account         string     `json:"account_id" validate:"omitempty"`
}
