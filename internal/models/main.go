package models

import (
	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"gorm.io/gorm"
)

type Models struct {
	Users interface {
		Create(data *dto.CreateUserDto) (*User, error)
		GetUserByEmail(email string) (*User, error)
		Save(user *User) error
	}
	Code interface {
		Create(code *Code) (*Code, error)
		GetByCode(code string) (*Code, error)
		Delete(code *Code) error
		GetByUserID(userID string) (*Code, error)
		GetByUserIDAndType(userID string, codeType utils.CodeType) (*Code, error)
		GetUnusedByUserIDAndType(userID string, codeType utils.CodeType) (*Code, error)
		GetByCodeAndType(code string, codeType utils.CodeType) (*Code, error)
		Save(code *Code) error
	}
}

func NewModels(db *gorm.DB) *Models {
	return &Models{
		Users: NewUserModel(db),
		Code:  NewCodeModel(db),
	}
}
