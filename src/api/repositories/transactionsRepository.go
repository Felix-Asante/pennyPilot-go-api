package repositories

import (
	"database/sql/driver"
	"encoding/json"
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
	AccountId   *string         `gorm:"column:account_id" json:"account,omitempty"`
	GoalId      *string         `gorm:"column:goal_id" json:"goal,omitempty"`
	Amount      float64         `gorm:"column:amount;default:0;not null" json:"amount"`
	Type        TransactionType `gorm:"type:transaction_type;column:type;not null"  json:"type"`
	Date        *time.Time      `gorm:"column:date;not null" json:"date"`
	Description string          `gorm:"column:description" json:"description"`
	CreatedAt   *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time      `gorm:"column:updated_at" json:"updated_at"`

	Account *Accounts `gorm:"foreignKey:AccountId" json:"-"`
	Goal    *Goals    `gorm:"foreignKey:GoalId" json:"-"`
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

func (t Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		Account *Accounts `json:"account,omitempty"`
		Goal    *Goals    `json:"goal,omitempty"`
		*Alias
	}{
		Account: t.Account,
		Goal:    t.Goal,
		Alias:   (*Alias)(&t),
	})
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

func (t *TransactionsRepository) FindAllByUserId(userId string, page int, pageSize int) (PaginationResult, error) {
	var transactions []Transaction
	query := t.db.Where("user_id = ?", userId).Find(&transactions).Preload("Account").Preload("Goal").Order("created_at desc")

	return Paginate(query, page, pageSize, &transactions)
}

func (t *TransactionsRepository) GetTransactionsByAccount(data GetAccountTransactions) (PaginationResult, error) {
	var transactions []Transaction

	query := t.db.
		Select("DATE_TRUNC('month', created_at) as month, *").
		Where("account_id = ?", data.AccountId).
		Group("DATE_TRUNC('month', created_at), id").
		Preload("Account").
		Order("DATE_TRUNC('month', created_at) desc")

	if !data.StartDate.IsZero() && !data.EndDate.IsZero() {
		query = query.Where("date BETWEEN ? AND ?", data.StartDate, data.EndDate)
	}

	if data.Type != "" {
		query = query.Where("type = ?", TransactionType(data.Type))
	}

	return Paginate(query, data.Page, data.PageSize, &transactions)
}
