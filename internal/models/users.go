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
	Codes        []Code    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{DB: db}
}

func (um *UserModel) Create(data *dto.CreateUserDto) (*User, error) {
	user := User{
		Email:        data.Email,
		FullName:     data.FullName,
		PasswordHash: data.Password,
		Currency:     data.Currency,
		ID:           uuid.New().String(),
	}

	err := um.DB.Create(&user).Error
	return &user, err
}

func (um *UserModel) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := um.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
