package repositories

import (
	"time"

	"gorm.io/gorm"
)

type Accounts struct {
	ID              string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	UserID          string     `gorm:"column:user_id;not null"`
	Name            string     `gorm:"column:name;unique;not null;index"`
	CurrentBalance  string     `gorm:"column:current_balance;not null"`
	TargetBalance   string     `gorm:"column:target_balance;not null"`
	AllocationPoint string     `gorm:"column:allocation_point;not null"`
	CreatedAt       *time.Time `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`
}

type AccountsRepository struct {
	db *gorm.DB
}

func (u *Accounts) TableName() string {
	return "accounts"
}

func NewAccountsRepository(db *gorm.DB) *AccountsRepository {

	return &AccountsRepository{db}
}
