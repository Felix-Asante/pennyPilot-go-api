package repositories

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Accounts struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id" json:"id"`
	UserID          string     `gorm:"column:user_id;not null" json:"-"`
	Name            string     `gorm:"column:name;not null;index" json:"name"`
	CurrentBalance  float64    `gorm:"column:current_balance;default:0" json:"current_balance"`
	TargetBalance   float64    `gorm:"column:target_balance;not null" json:"target_balance"`
	AllocationPoint float64    `gorm:"column:allocation_point;not null" json:"allocation_point"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalAllocation float64    `gorm:"column:total_allocation;default:0"`

	Goals        []Goals       `gorm:"columns:goals;foreignKey:account_id;reference:id"`
	Transactions []Transaction `gorm:"columns:transactions;foreignKey:account_id;reference:id"`
	Expenses     []Expenses    `gorm:"columns:expenses;foreignKey:account_id;reference:id"`
}

type AccountsRepository struct {
	db *gorm.DB
}

type CreateAccountDto struct {
	UserID          string  `json:"user_id"`
	Name            string  `json:"name"`
	TargetBalance   float64 `json:"target_balance"`
	AllocationPoint float64 `json:"allocation_point"`
}

type NewAccountResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	CurrentBalance  float64    `json:"current_balance"`
	TargetBalance   float64    `json:"target_balance"`
	AllocationPoint float64    `json:"allocation_point"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

func (u *Accounts) TableName() string {
	return "accounts"
}

func (u *AccountsRepository) Create(data CreateAccountDto) (*NewAccountResponse, error) {
	account := Accounts{
		UserID:          data.UserID,
		Name:            data.Name,
		TargetBalance:   data.TargetBalance,
		AllocationPoint: data.AllocationPoint,
		CurrentBalance:  0.00,
	}

	error := u.db.Create(&account).Error

	newAccount := NewAccountResponse{
		ID:              account.ID.String(),
		Name:            account.Name,
		CurrentBalance:  account.CurrentBalance,
		TargetBalance:   account.TargetBalance,
		AllocationPoint: account.AllocationPoint,
		CreatedAt:       account.CreatedAt,
		UpdatedAt:       account.UpdatedAt,
	}

	return &newAccount, error
}

func (u *AccountsRepository) FindByNameAndUserID(name string, userID string) (*Accounts, error) {

	var existingAccount Accounts

	error := u.db.Where("name = ? AND user_id = ?", name, userID).Find(&existingAccount).Error

	return &existingAccount, error
}

func (u *AccountsRepository) FindByIDAndUserID(id string, userID string) (*Accounts, error) {

	var existingAccount Accounts

	error := u.db.Where("id = ? AND user_id = ?", id, userID).Find(&existingAccount).Error

	if error != nil {
		return nil, error
	}
	return &existingAccount, error
}

func (u *AccountsRepository) Save(account *Accounts) (*Accounts, error) {
	error := u.db.Save(account).Error

	return account, error
}

func (u *AccountsRepository) Remove(id string) (bool, error) {
	var accounts Accounts

	error := u.db.Where("id", id).Delete(&accounts).Error

	if error != nil {
		return false, error
	}

	return true, nil
}

func (u *AccountsRepository) FindAllByUserID(userId string, queries AccountQueries) (PaginationResult, error) {
	var accounts []Accounts

	query := u.db.Where("user_id = ?", userId)

	if queries.Sort != "" {
		parts := strings.Split(queries.Sort, ":")

		if len(parts) == 2 {
			field := parts[0]
			order := parts[1]
			query = query.Order(field + " " + order)
		}
	}

	if queries.Query != "" {
		query = query.Where("name ILIKE ?", "%"+queries.Query+"%")
	}

	query = query.Find(&accounts)

	return Paginate(query, queries.Page, queries.Limit, &accounts)
}

func (u *AccountsRepository) FindTotalSavings(userId string) (float64, error) {
	var totalSavings float64
	err := u.db.Model(&Accounts{}).
		Select("COALESCE(SUM(current_balance + total_allocation), 0) as total_savings").
		Where("user_id = ?", userId).
		Row().
		Scan(&totalSavings)

	if err != nil {
		return 0, err
	}
	return totalSavings, nil
}

func NewAccountsRepository(db *gorm.DB) *AccountsRepository {

	return &AccountsRepository{db}
}

func (u *Accounts) BeforeCreate(tx *gorm.DB) error {

	u.ID = uuid.New()
	return nil
}
