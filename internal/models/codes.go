package models

import (
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"gorm.io/gorm"
)

type Code struct {
	gorm.Model
	UserID    string         `gorm:"column:user_id;index"`
	Code      string         `gorm:"column:code;unique;not null;index"`
	ExpiresAt *time.Time     `gorm:"column:expires_at"`
	Used      bool           `gorm:"column:used"`
	Type      utils.CodeType `gorm:"column:type;index;type:code_type"`
}

type CodeModel struct {
	DB *gorm.DB
}

func NewCodeModel(db *gorm.DB) *CodeModel {
	return &CodeModel{DB: db}
}

func (cm *CodeModel) Create(code *Code, tx *gorm.DB) (*Code, error) {
	if tx != nil {
		err := tx.Create(&code).Error
		return code, err
	}
	err := cm.DB.Create(&code).Error
	return code, err
}

func (cm *CodeModel) GetByCode(code string) (*Code, error) {
	var existingCode Code
	if err := cm.DB.Where("code = ?", code).First(&existingCode).Error; err != nil {
		return nil, err
	}
	return &existingCode, nil
}

func (cm *CodeModel) Delete(code *Code, tx *gorm.DB) error {
	if tx != nil {
		return tx.Delete(code).Error
	}
	return cm.DB.Delete(code).Error
}

func (cm *CodeModel) GetByUserID(userID string) (*Code, error) {
	var existingCode Code
	if err := cm.DB.Where("user_id = ?", userID).First(&existingCode).Error; err != nil {
		return nil, err
	}
	return &existingCode, nil
}

func (cm *CodeModel) GetByCodeAndType(code string, codeType utils.CodeType) (*Code, error) {
	var existingCode Code
	if err := cm.DB.Where("code = ? AND type = ?", code, codeType).First(&existingCode).Error; err != nil {
		return nil, err
	}
	return &existingCode, nil
}

func (cm *CodeModel) GetByUserIDAndType(userID string, codeType utils.CodeType) (*Code, error) {
	var existingCode Code
	if err := cm.DB.Where("user_id = ? AND type = ?", userID, codeType).First(&existingCode).Error; err != nil {
		return nil, err
	}
	return &existingCode, nil
}

func (cm *CodeModel) GetUnusedByUserIDAndType(userID string, codeType utils.CodeType) (*Code, error) {
	var existingCode Code
	if err := cm.DB.Where("user_id = ? AND type = ? AND used = ?", userID, codeType, false).First(&existingCode).Error; err != nil {
		return nil, err
	}
	return &existingCode, nil
}

func (cm *CodeModel) Save(code *Code, tx *gorm.DB) error {
	if tx != nil {
		return tx.Save(code).Error
	}
	return cm.DB.Save(code).Error
}
