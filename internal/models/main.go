package models

import (
	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"gorm.io/gorm"
)

type Models struct {
	Users interface {
		Create(*dto.CreateUserDto) (*User, error)
		GetUserByEmail(string) (*User, error)
	}
}

func NewModels(db *gorm.DB) *Models {
	return &Models{
		Users: NewUserModel(db),
	}
}
