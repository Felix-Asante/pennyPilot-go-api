package dto

import (
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
)

type CreateIncomeDto struct {
	Amount       float64               `json:"amount" validate:"required" errormgs:"Amount is required"`
	Category     *string               `json:"category" validate:"required" errormgs:"Category is required"`
	DateRecieved time.Time             `json:"date_recieved" validate:"required" errormgs:"Date recieved is required"`
	Type         utils.IncomeType      `json:"type" validate:"required,oneof=salary freelance investment other" errormgs:"Type is required"`
	Frequency    utils.IncomeFrequency `json:"frequency" validate:"required,oneof=weekly biweekly monthly yearly one-time" errormgs:"Frequency is required"`
}

type UpdateIncomeDto struct {
	Amount       *float64               `json:"amount" validate:"omitempty" errormgs:"Amount is required"`
	Category     *string                `json:"category" validate:"omitempty" errormgs:"Category is required"`
	DateRecieved *time.Time             `json:"date_recieved" validate:"omitempty,gte" errormgs:"Date recieved is required"`
	Type         *utils.IncomeType      `json:"type" validate:"omitempty,oneof=salary freelance investment other" errormgs:"Type is required"`
	Frequency    *utils.IncomeFrequency `json:"frequency" validate:"omitempty,oneof=weekly biweekly monthly yearly one-time" errormgs:"Frequency is required"`
}

type TransferIncome struct {
	Amount   float64  `json:"amount" validate:"required,min=1"`
	Accounts []string `json:"accounts" validate:"required"`
}
