package models

import (
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey;column:id;index"`
	Email        string    `gorm:"column:email;unique;index"`
	FullName     string    `gorm:"column:full_name"`
	PasswordHash string    `gorm:"column:password_hash"`
	Currency     string    `gorm:"column:currency"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	TotalIncome  float64   `gorm:"-:all"`

	Codes         []Code         `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Incomes       []Income       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Accounts      []Account      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	IncomeBalance *IncomeBalance `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{DB: db}
}

func (um *UserModel) Create(data *dto.CreateUserDto, tx *gorm.DB) (*User, error) {
	user := User{
		Email:        data.Email,
		FullName:     data.FullName,
		PasswordHash: data.Password,
		Currency:     data.Currency,
		ID:           uuid.New().String(),
	}

	db := getTxDB(um.DB, tx)

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (um *UserModel) GetUserByEmail(email string, tx *gorm.DB) (*User, error) {
	var user User
	db := getTxDB(um.DB, tx)

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	var total float64
	err := db.Model(&Income{}).
		Select("COALESCE(SUM(amount),0)").
		Where("user_id = ?", user.ID).
		Scan(&total).Error

	if err != nil {
		return &user, err
	}

	user.TotalIncome = total

	return &user, nil
}

func (um *UserModel) Save(user *User, tx *gorm.DB) error {
	db := getTxDB(um.DB, tx)
	return db.Save(user).Error
}
