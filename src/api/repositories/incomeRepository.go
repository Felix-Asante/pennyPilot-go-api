package repositories

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomeType string
type IncomeFrequency string

const (
	Salary     IncomeType = "salary"
	Investment IncomeType = "investment"
	Freelance  IncomeType = "freelance"
	Others     IncomeType = "others"
)
const (
	Weekly   IncomeFrequency = "weekly"
	Monthly  IncomeFrequency = "monthly"
	Yearly   IncomeFrequency = "yearly"
	BiWeekly IncomeFrequency = "bi-weekly"
	OneTime  IncomeFrequency = "one-time"
)

type Incomes struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id" json:"id"`
	UserId       string          `gorm:"type:uuid;column:user_id;not null" json:"-"`
	Amount       float64         `gorm:"column:amount;default:0" json:"amount"`
	DateReceived *time.Time      `gorm:"column:date_received;not null" json:"date_received"`
	Type         IncomeType      `gorm:"type:income_type;column:type;not null"  json:"type"`
	Frequency    IncomeFrequency `gorm:"type:income_frequency;column:frequency;not null"  json:"frequency"`
	CreatedAt    *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (a *IncomeType) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid value for Income type")
	}

	*a = IncomeType(strValue)
	return nil
}

func (a IncomeType) Value() (driver.Value, error) {
	return string(a), nil
}

func (a *IncomeFrequency) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid value for Frequency type")
	}

	*a = IncomeFrequency(strValue)
	return nil
}

func (a IncomeFrequency) Value() (driver.Value, error) {
	return string(a), nil
}

func (u *Incomes) TableName() string {
	return "incomes"
}

type IncomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) *IncomeRepository {

	return &IncomeRepository{db}
}

func (repo *IncomeRepository) Create(data CreateIncomeDto) (*Incomes, error) {
	newIncome := Incomes{
		UserId:       data.User,
		Type:         IncomeType(data.Type),
		Frequency:    IncomeFrequency(data.Frequency),
		DateReceived: data.DateReceived,
		Amount:       data.Amount,
	}
	error := repo.db.Create(&newIncome).Error

	return &newIncome, error
}

func (u *Incomes) BeforeCreate(tx *gorm.DB) error {

	u.ID = uuid.New()
	return nil
}
