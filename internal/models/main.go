package models

import "gorm.io/gorm"

type Models struct {
	Users interface {
		Create(*User) error
		GetUserByEmail(string) (*User, error)
	}
}

func NewModels(db *gorm.DB) *Models {
	return &Models{
		Users: NewUserModel(db),
	}
}
