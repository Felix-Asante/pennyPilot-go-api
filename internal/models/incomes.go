package models

import (
	"context"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Income struct {
	ID           uuid.UUID             `gorm:"type:uuid;primaryKey;column:id;index"`
	UserID       string                `gorm:"column:user_id;index"`
	Amount       float64               `gorm:"column:amount;not null"`
	Category     *string               `gorm:"column:category;"`
	DateRecieved time.Time             `gorm:"column:date_recieved;not null"`
	Type         utils.IncomeType      `gorm:"column:type;type:income_type;not null"`
	Frequency    utils.IncomeFrequency `gorm:"column:frequency;type:income_frequency;not null"`
	CreatedAt    time.Time             `gorm:"autoCreateTime"`
	UpdatedAt    time.Time             `gorm:"autoUpdateTime"`
}

type IncomeModel struct {
	DB *gorm.DB
}

func NewIncomeModel(db *gorm.DB) *IncomeModel {
	return &IncomeModel{DB: db}
}

func (im *IncomeModel) Create(income *Income, tx *gorm.DB) (*Income, error) {
	db := getTxDB(im.DB, tx)

	if err := db.Create(income).Error; err != nil {
		return nil, err
	}

	return income, nil
}

func (im *IncomeModel) GetAllByUserID(userID string, tx *gorm.DB) ([]*Income, error) {
	db := getTxDB(im.DB, tx)

	var income []*Income
	if err := db.Where("user_id = ?", userID).Find(&income).Error; err != nil {
		return nil, err
	}

	return income, nil
}

func (im *IncomeModel) GetByID(id string, tx *gorm.DB) (*Income, error) {
	db := getTxDB(im.DB, tx)

	var income Income
	if err := db.Where("id = ?", id).First(&income).Error; err != nil {
		return nil, err
	}

	return &income, nil
}

func (im *IncomeModel) Save(income *Income, tx *gorm.DB) error {
	db := getTxDB(im.DB, tx)

	if err := db.Save(income).Error; err != nil {
		return err
	}

	return nil
}

func (im *IncomeModel) GetUserTotalIncome(ctx context.Context, userId string, tx *gorm.DB) (float64, error) {
	db := getTxDB(im.DB, tx)
	var total float64

	err := db.Model(&Income{}).
		Select("COALESCE(SUM(amount),0)").
		Where("user_id = ?", userId).
		Scan(&total).Error

	return total, err
}
