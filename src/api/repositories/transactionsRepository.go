package repositories

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
	Allocation TransactionType = "allocation"
	Expense    TransactionType = "expense"
	Interest   TransactionType = "interest"
	Penalty    TransactionType = "penalty"
	Refund     TransactionType = "refund"
)

type Transaction struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id" json:"id"`
	UserId      string          `gorm:"type:uuid;column:user_id;not null" json:"user"`
	AccountId   *string         `gorm:"column:account_id" json:"account"`
	GoalId      *string         `gorm:"column:goal_id" json:"goal"`
	Amount      float64         `gorm:"column:amount;default:0;not null" json:"amount"`
	Type        TransactionType `gorm:"type:transaction_type;column:type;not null"  json:"type"`
	Date        *time.Time      `gorm:"column:date;not null" json:"date"`
	Description string          `gorm:"column:description" json:"description"`
	CreatedAt   *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (a *TransactionType) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid value for transaction type")
	}

	*a = TransactionType(strValue)
	return nil
}

func (a TransactionType) Value() (driver.Value, error) {
	return string(a), nil
}

func (t *Transaction) TableName() string {
	return "transactions"
}

type TransactionsRepository struct {
	db *gorm.DB
}

func NewTransactionsRepository(db *gorm.DB) *TransactionsRepository {
	return &TransactionsRepository{db}
}

func (t *TransactionsRepository) Create(data CreateTransactionDto) (*Transaction, error) {
	newTransaction := Transaction{
		UserId:      data.User,
		AccountId:   data.Account,
		GoalId:      data.Goal,
		Amount:      data.Amount,
		Type:        TransactionType(data.Type),
		Date:        data.Date,
		Description: data.Description,
	}
	error := t.db.Create(&newTransaction).Error

	return &newTransaction, error
}
