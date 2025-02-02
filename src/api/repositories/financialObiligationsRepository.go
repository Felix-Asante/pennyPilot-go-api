package repositories

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinancialObligationType string
type FinancialObligationRepaymentType string

const (
	Debt FinancialObligationType = "debt"
	Loan FinancialObligationType = "loan"
)
const (
	PaidWeekly  FinancialObligationRepaymentType = "weekly"
	PaidMonthly FinancialObligationRepaymentType = "monthly"
	PaidYearly  FinancialObligationRepaymentType = "yearly"
	PaidDaily   FinancialObligationRepaymentType = "daily"
)

type FinancialObligations struct {
	ID               uuid.UUID                        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id" json:"id"`
	UserID           string                           `gorm:"column:user_id;not null" json:"user"`
	Type             FinancialObligationType          `gorm:"type:financial_obligation_type;column:financial_obligation_type;not null"  json:"financial_obligation_type"`
	TotalAmount      float64                          `gorm:"column:total_amount;not null" json:"total_amount"`
	CounterpartyName string                           `gorm:"column:counterparty_name;not null" json:"counterparty_name"`
	RemainingAmount  float64                          `gorm:"column:remaining_amount;not null" json:"remaining_amount"`
	RepaymentType    FinancialObligationRepaymentType `gorm:"type:financial_obligation_repayment_type;column:financial_obligation_repayment_type;not null"  json:"financial_obligation_repayment_type"`
	InterestRate     float64                          `gorm:"column:interest_rate;" json:"interest_rate"`
	NextDueDate      *time.Time                       `gorm:"column:next_due_date;" json:"next_due_date"`
	CreatedAt        *time.Time                       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        *time.Time                       `gorm:"column:updated_at" json:"updated_at"`
}

type FinancialObligationsRepository struct {
	db *gorm.DB
}

type CreateFinancialObligation struct {
	Type             FinancialObligationType          `json:"type" validate:"required,oneof=loan debt"`
	TotalAmount      float64                          `json:"total_amount" validate:"required,gt=0"`
	CounterpartyName string                           `json:"counterparty_name" validate:"required,min=2"`
	RemainingAmount  float64                          `json:"remaining_amount" validate:"omitempty,min=0"`
	RepaymentType    FinancialObligationRepaymentType `json:"repayment_type" validate:"required,oneof=weekly monthly yearly daily"`
	InterestRate     float64                          `json:"interest_rate" validate:"omitempty,gt=0"`
	NextDueDate      *time.Time                       `json:"next_due_date" validate:"omitempty"`
	UserId           string                           `json:"user_id"`
}

func (a *FinancialObligationType) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid value for financial obligation type")
	}

	*a = FinancialObligationType(strValue)
	return nil
}

func (a FinancialObligationType) Value() (driver.Value, error) {
	return string(a), nil
}

func (a *FinancialObligationRepaymentType) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid value for Financial obligation repayment type")
	}

	*a = FinancialObligationRepaymentType(strValue)
	return nil
}

func (a FinancialObligationRepaymentType) Value() (driver.Value, error) {
	return string(a), nil
}

func (u *FinancialObligations) TableName() string {
	return "financial_obligations"
}

func NewFinancialObligationsRepository(db *gorm.DB) *FinancialObligationsRepository {

	return &FinancialObligationsRepository{db}
}

func (u *FinancialObligationsRepository) Create(dto CreateFinancialObligation) (*FinancialObligations, error) {
	obligation := FinancialObligations{
		Type:             FinancialObligationType(dto.Type),
		TotalAmount:      dto.TotalAmount,
		CounterpartyName: dto.CounterpartyName,
		RemainingAmount:  dto.RemainingAmount,
		RepaymentType:    FinancialObligationRepaymentType(dto.RepaymentType),
		InterestRate:     dto.InterestRate,
		NextDueDate:      dto.NextDueDate,
		UserID:           dto.UserId,
	}
	error := u.db.Create(&obligation).Error
	return &obligation, error
}

func (u *FinancialObligationsRepository) FindByID(id string) (*FinancialObligations, error) {
	var obligation FinancialObligations
	error := u.db.Where("id = ?", id).Find(&obligation).Error
	return &obligation, error
}

func (u *FinancialObligationsRepository) FindByIDAndUserID(id string, userId string) (*FinancialObligations, error) {
	var obligation FinancialObligations
	error := u.db.Where("id = ? and user_id = ?", id, userId).Find(&obligation).Error
	return &obligation, error
}

func (u *FinancialObligationsRepository) Remove(id string) (bool, error) {
	var obligation FinancialObligations
	error := u.db.Where("id = ?", id).Delete(&obligation).Error
	return error == nil, error
}

func (u *FinancialObligationsRepository) FindByUserID(id string) ([]FinancialObligations, error) {
	var obligations []FinancialObligations
	error := u.db.Where("user_id = ?", id).Find(&obligations).Error
	return obligations, error
}
