package utils

import (
	"database/sql/driver"
	"fmt"
)

type CodeType string
type IncomeType string
type IncomeFrequency string
type AllocationStrategy string

const (
	CodeTypeForgotPassword CodeType = "forgot_password"
	CodeTypeVerifyEmail    CodeType = "verify_email"
)

const (
	IncomeTypeSalary     IncomeType = "salary"
	IncomeTypeFreelance  IncomeType = "freelance"
	IncomeTypeInvestment IncomeType = "investment"
	IncomeTypeOther      IncomeType = "other"
)

const (
	IncomeFrequencyWeekly   IncomeFrequency = "weekly"
	IncomeFrequencyBiweekly IncomeFrequency = "biweekly"
	IncomeFrequencyMonthly  IncomeFrequency = "monthly"
	IncomeFrequencyYearly   IncomeFrequency = "yearly"
	IncomeFrequencyOneTime  IncomeFrequency = "one-time"
)

const (
	AllocationStrategyFixedAmount AllocationStrategy = "fixed_amount"
	AllocationStrategyPercentage  AllocationStrategy = "percentage"
)

func (p *CodeType) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = CodeType(v)
	case string:
		*p = CodeType(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (p CodeType) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *IncomeType) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = IncomeType(v)
	case string:
		*p = IncomeType(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (p IncomeType) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *IncomeFrequency) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = IncomeFrequency(v)
	case string:
		*p = IncomeFrequency(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (p IncomeFrequency) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *AllocationStrategy) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*p = AllocationStrategy(v)
	case string:
		*p = AllocationStrategy(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (p AllocationStrategy) Value() (driver.Value, error) {
	return string(p), nil
}
