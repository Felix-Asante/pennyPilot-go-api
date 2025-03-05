package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expenses struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id" json:"id"`
	Account     *string  `gorm:"column:account_id" json:"account"`
	Amount      float64    `gorm:"column:amount;not null" json:"amount"`
	Description string     `gorm:"column:description;not null" json:"description"`
	Date        *time.Time `gorm:"column:date;not null" json:"date"`
	Category   int  `gorm:"column:category_id;not null" json:"-"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`

	CategoryRef ExpenseCategory `gorm:"foreignKey:Category" json:"category"`
}

type ExpenseCategory struct {
	gorm.Model
	Name     string     `gorm:"column:name;not null" json:"name"`
	Icon     string     `gorm:"column:icon;not null" json:"icon"`
	UserID   string     `gorm:"column:user_id;not null" json:"-"`
	Expenses []Expenses `gorm:"columns:expenses;foreignKey:category_id;reference:id"`
}

type ExpensesRepository struct {
	db *gorm.DB
}
type ExpenseCategoryRepository struct {
	db *gorm.DB
}

func (expCategory *ExpenseCategoryRepository) Create(data CreateExpenseCategoryDto) (ExpenseCategory, error) {
	newCategory := ExpenseCategory{
		UserID: data.User,
		Name:   data.Name,
		Icon:   data.Icon,
	}
	error := expCategory.db.Create(&newCategory).Error

	return newCategory, error
}

func (expCategory *ExpenseCategoryRepository) FindAllByUserID(id string) ([]ExpenseCategory, error) {
	var categories []ExpenseCategory
	error := expCategory.db.Where("user_id = ?", id).Find(&categories).Error

	return categories,error
}

func (expRepo *ExpensesRepository) Delete(id string) error {
	error := expRepo.db.Where("id =?", id).Delete(&Expenses{}).Error
	return error
}
func (expRepo *ExpensesRepository) GetByID(id string) (*Expenses, error) {
	var expense Expenses
	error := expRepo.db.Where("id =?", id).First(&expense).Error
	return &expense, error
}

func (expCateg *ExpenseCategoryRepository) GetByID(id int) (*ExpenseCategory, error) {
	var category ExpenseCategory
	error := expCateg.db.Where("id =?", id).First(&category).Error
	return &category, error
}

func (expCateg *ExpenseCategoryRepository) Delete(id int) error {
	error := expCateg.db.Where("id =?", id).Delete(&ExpenseCategory{}).Error
	return error
}

func (expRepo *ExpensesRepository) GetByAccount(data GetAccountExpensesDto) ([]Expenses, error) {
	var expenses []Expenses
	query := expRepo.db.Where("account_id = ? AND EXTRACT(YEAR FROM date) = ?", data.Account, data.Year).Preload("CategoryRef")

	if !data.StartDate.IsZero() && !data.EndDate.IsZero() {
		query = query.Where("date BETWEEN ? AND ?", data.StartDate, data.EndDate)
	}

	error := query.Find(&expenses).Error
	return expenses, error
}

func (expRepo *ExpensesRepository) GetByCategory(data GetExpenseByCategoryDto) (PaginationResult, error) {
	var expenses []Expenses
	query := expRepo.db.Where("category_id = ?", data.Category).Find(&expenses)

	return Paginate(query, data.Pagination.Page, data.Pagination.Limit, &expenses)
}

func (expRepo *ExpensesRepository) GetUserTotal(user string) (float64, error) {
	var totalExpense float64
	err := expRepo.db.Model(&Expenses{}).Select("COALESCE(SUM(expense.amount), 0)").Joins("JOIN expense_category ON expense.category_id = expense_category.id").
	Where("expense_category.user_id =?", user).Row().Scan(&totalExpense)
	if err != nil {
		return 0, err
	}
	return totalExpense, nil
}


func (expRepo *ExpensesRepository) GetTotalByAccount(data GetAccountExpensesDto) (float64, error) {
	var totalExpense float64
	err := expRepo.db.Model(&Expenses{}).Select("COALESCE(SUM(amount), 0)").
	Where("account_id =?", data.Account).Row().Scan(&totalExpense)
	if err != nil {
		return 0, err
	}
	return totalExpense, nil
}

func (expRepo *ExpensesRepository) Get(data GetExpenseDto) (PaginationResult, error) {
	var expenses []Expenses

	query := expRepo.db.Model(&Expenses{}).
		Joins("JOIN expense_category ON expense.category_id = expense_category.id").
		Where("expense_category.user_id =? AND EXTRACT(YEAR FROM date) =?", data.User, data.Year)

	if !data.StartDate.IsZero() && !data.EndDate.IsZero() {
		query = query.Where("date BETWEEN ? AND ?", data.StartDate, data.EndDate)
	}

	return Paginate(query, data.Pagination.Page, data.Pagination.Limit, &expenses)
}

func (expRepo *ExpensesRepository) Create(data CreateExpenseDto) (*Expenses, error) {
	newExpense := Expenses{
		Account:     &data.Account,
		Amount:      data.Amount,
		Description: data.Description,
		Date:        data.Date,
		Category:    data.Category,
	}
	error := expRepo.db.Create(&newExpense).Error

	return &newExpense, error
}

func NewExpenseCategoryRepository(db *gorm.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{
		db: db,
	}
}

func NewExpensesRepository(db *gorm.DB) *ExpensesRepository {
	return &ExpensesRepository{
		db: db,
	}
}

func NewExpenseRepository(db *gorm.DB) *ExpensesRepository {
	return &ExpensesRepository{
		db: db,
	}
}

func (u *Expenses) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

func (u *Expenses) TableName() string {
	return "expense"
}

func (u *ExpenseCategory) TableName() string {
	return "expense_category"
}
