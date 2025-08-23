package models

import (
	"context"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/dto"
	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/google/uuid"
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
	Account interface {
		Create(ctx context.Context, account *Account, tx *gorm.DB) error
		Save(ctx context.Context, account *Account, tx *gorm.DB) error
		Delete(ctx context.Context, account *Account, tx *gorm.DB) error
		GetAllByUserID(ctx context.Context, userID string, tx *gorm.DB) ([]*Account, error)
		GetByID(ctx context.Context, id string, tx *gorm.DB) (*Account, error)
		GetByIDAndUserID(ctx context.Context, id string, userID string, tx *gorm.DB) (*Account, error)
		GetByNameAndUserID(ctx context.Context, name string, userID string, tx *gorm.DB) (*Account, error)
	}
	Envelope interface {
		Create(ctx context.Context, envelope *Envelope, tx *gorm.DB) error
		Save(ctx context.Context, envelope *Envelope, tx *gorm.DB) error
		Delete(ctx context.Context, envelope *Envelope, tx *gorm.DB) error
		GetAllByAccountID(ctx context.Context, accountID uuid.UUID, tx *gorm.DB) ([]*Envelope, error)
		GetByID(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*Envelope, error)
		GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID, tx *gorm.DB) (*Envelope, error)
	}
	AllocationRule interface {
		Create(ctx context.Context, allocationRule *AllocationRule, tx *gorm.DB) error
		Save(ctx context.Context, allocationRule *AllocationRule, tx *gorm.DB) error
		Delete(ctx context.Context, allocationRule *AllocationRule, tx *gorm.DB) error
		GetByTargetID(ctx context.Context, targetID uuid.UUID, tx *gorm.DB) ([]*AllocationRule, error)
		GetByID(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*AllocationRule, error)
		GetByIDAndTargetID(ctx context.Context, id uuid.UUID, targetID uuid.UUID, tx *gorm.DB) (*AllocationRule, error)
	}
}

func NewModels(db *gorm.DB) *Models {
	return &Models{
		Users:          NewUserModel(db),
		Code:           NewCodeModel(db),
		Income:         NewIncomeModel(db),
		Account:        NewAccountModel(db),
		Envelope:       NewEnvelopeModel(db),
		AllocationRule: NewAllocationRuleModel(db),
	}
}
