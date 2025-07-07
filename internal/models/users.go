package models

import "gorm.io/gorm"

type User struct {
	ID           string `gorm:"type:uuid;primaryKey;column:id;index"`
	Email        string `gorm:"column:email;unique;index"`
	FullName     string `gorm:"column:full_name"`
	PasswordHash string `gorm:"column:password_hash"`
	Currency     string `gorm:"column:currency"`
	CreatedAt    int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt    int64  `gorm:"autoUpdateTime:milli"`
}

type UserModel struct {
	DB *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{DB: db}
}

func (um *UserModel) Create(user *User) error {
	return um.DB.Create(user).Error
}

func (um *UserModel) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := um.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
