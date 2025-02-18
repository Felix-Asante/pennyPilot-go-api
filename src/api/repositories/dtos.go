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

type CreateFinancialObligationDto struct {
	Type             FinancialObligationType          `json:"type" validate:"required,oneof=loan debt"`
	TotalAmount      float64                          `json:"total_amount" validate:"required,gt=0"`
	CounterpartyName string                           `json:"counterparty_name" validate:"required,min=2"`
	RemainingAmount  float64                          `json:"remaining_amount" validate:"omitempty,min=0"`
	RepaymentType    FinancialObligationRepaymentType `json:"repayment_type" validate:"required,oneof=weekly monthly yearly daily"`
	InterestRate     float64                          `json:"interest_rate" validate:"omitempty,gt=0"`
	NextDueDate      *time.Time                       `json:"next_due_date" validate:"omitempty"`
}

type CreateTransactionDto struct {
	Account     *string         `json:"account_id" validate:"omitempty"`
	Goal        *string         `json:"goal_id" validate:"omitempty"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	User        string          `json:"user_id" validate:"required"`
	Type        TransactionType `json:"type" validate:"required,oneof=deposit withdrawal transfer allocation expense"`
	Date        *time.Time      `json:"date" validate:"required"`
	Description string          `json:"description" validate:"omitempty,min=2"`
}

type AccountQueries struct {
	Page  int    `json:"page" validate:"omitempty"`
	Limit int    `json:"limit" validate:"omitempty"`
	Query string `json:"query" validate:"omitempty"`
	Sort  string `json:"sort" validate:"omitempty"`
}

type GetAccountTransactions struct {
	AccountId string `json:"account_id" validate:"required"`
	Page      int    `json:"page" validate:"omitempty"`
	PageSize  int    `json:"page_size" validate:"omitempty"`
	UserId    string `json:"user_id" validate:"required"`
}
