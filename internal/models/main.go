package models

import (
	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"gorm.io/gorm"
)

type Models struct {
	Users interface {
		Create(data *dto.CreateUserDto, tx *gorm.DB) (*User, error)
		GetUserByEmail(email string, tx *gorm.DB) (*User, error)
		Save(user *User, tx *gorm.DB) error
	}
	Code interface {
		Create(code *Code, tx *gorm.DB) (*Code, error)
		GetByCode(code string) (*Code, error)
		Delete(code *Code, tx *gorm.DB) error
		GetByUserID(userID string) (*Code, error)
		GetByUserIDAndType(userID string, codeType utils.CodeType) (*Code, error)
		GetUnusedByUserIDAndType(userID string, codeType utils.CodeType) (*Code, error)
		GetByCodeAndType(code string, codeType utils.CodeType) (*Code, error)
		Save(code *Code, tx *gorm.DB) error
	}
	Income interface {
		Create(income *Income, tx *gorm.DB) (*Income, error)
		GetAllByUserID(userID string, tx *gorm.DB) ([]*Income, error)
		GetByID(id string, tx *gorm.DB) (*Income, error)
		Save(income *Income, tx *gorm.DB) error
	}
}

func NewModels(db *gorm.DB) *Models {
	return &Models{
		Users:  NewUserModel(db),
		Code:   NewCodeModel(db),
		Income: NewIncomeModel(db),
	}
}
