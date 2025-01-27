package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Goals struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id" json:"id"`
	Name            string     `gorm:"column:name;not null" json:"name"`
	TargetAmount    float64    `gorm:"column:target_amount;not null" json:"target_amount"`
	CurrentBalance  float64    `gorm:"column:current_balance;not null" json:"current_balance"`
	AllocationPoint float64    `gorm:"column:allocation_point;not null" json:"allocation_point"`
	DueDate         *time.Time `gorm:"column:due_date;not null" json:"due_date"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
	AccountId       string     `gorm:"column:account_id;not null" json:"account_id"`
}

type GoalsRepository struct {
	db *gorm.DB
}

type CreateNewGoalData struct {
	Name            string     `json:"name" validate:"required,min=2"`
	TargetAmount    float64    `json:"target_amount" validate:"required,min=0.00"`
	AllocationPoint float64    `json:"allocation_point" validate:"required,min=1,max=100"`
	DueDate         *time.Time `json:"due_date" validate:"required,min=1,max=100"`
	Account         string     `json:"account_id" validate:"required"`
	CurrentBalance  float64    `json:"crrent_balance" validate:"required,min=1,max=100"`
}

func (u *Goals) TableName() string {
	return "goals"
}

func NewGoalsRepository(db *gorm.DB) *GoalsRepository {

	return &GoalsRepository{db}
}

func (repo *GoalsRepository) Create(data CreateNewGoalData) (*Goals, error) {
	newGoals := Goals{
		Name:            data.Name,
		AccountId:       data.Account,
		TargetAmount:    data.TargetAmount,
		CurrentBalance:  data.CurrentBalance,
		AllocationPoint: data.AllocationPoint,
		DueDate:         data.DueDate,
	}
	error := repo.db.Create(&newGoals).Error

	return &newGoals, error
}

func (repo *GoalsRepository) FindAccountTotalAllocation(accountId string) (float64, error) {
	var totalAllocation float64

	error := repo.db.Model(&Goals{}).Where("account_id = ?", accountId).Select("SUM(allocation_point)").Scan(&totalAllocation).Error

	return totalAllocation, error
}

func (u *Goals) BeforeCreate(tx *gorm.DB) error {

	u.ID = uuid.New()
	return nil
}
