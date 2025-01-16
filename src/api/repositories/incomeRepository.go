package repositories

import (
	"time"

	"github.com/google/uuid"
)

type IncomeType string
type IncomeFrequency string

type Incomes struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id" json:"id"`
	UserId       string          `gorm:"type:uuid;column:user_id" json:"-"`
	Amount       float64         `gorm:"column:amount;default:0" json:"amount"`
	DateReceived *time.Time      `gorm:"column:date_received" json:"date_received"`
	Type         IncomeType      `sql:"type:ENUM('salary', 'investment', 'freelance','others')" gorm:"column:type" json:"type"`
	Frequency    IncomeFrequency `sql:"type:ENUM('one-time', 'weekly', 'bi-weekly','monthly','yearly')" gorm:"column:frequency" json:"frequency"`
	CreatedAt    *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (u *Incomes) TableName() string {
	return "incomes"
}
