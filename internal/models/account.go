package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;column:id;index"`
	Name      string         `gorm:"column:name;not null"`
	UserID    string         `gorm:"column:user_id;index"`
	Balance   float64        `gorm:"column:balance;not null;default:0.00"`
	Currency  string         `gorm:"column:currency;not null"`
	IsActive  bool           `gorm:"column:is_active;not null;default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Envelopes []Envelope `gorm:"foreignKey:AccountID;references:ID"`
}

type AccountModel struct {
	db *gorm.DB
}

func NewAccountModel(db *gorm.DB) *AccountModel {
	return &AccountModel{db: db}
}

func (am *AccountModel) Create(ctx context.Context, account *Account, tx *gorm.DB) error {
	db := getTxDB(am.db, tx)

	if err := db.WithContext(ctx).Create(account).Error; err != nil {
		return err
	}

	return nil
}

func (am *AccountModel) Save(ctx context.Context, account *Account, tx *gorm.DB) error {
	db := getTxDB(am.db, tx)

	if err := db.WithContext(ctx).Save(account).Error; err != nil {
		return err
	}

	return nil
}

func (am *AccountModel) Delete(ctx context.Context, account *Account, tx *gorm.DB) error {
	db := getTxDB(am.db, tx)

	if err := db.WithContext(ctx).Delete(account).Error; err != nil {
		return err
	}

	return nil
}

func (am *AccountModel) GetAllByUserID(ctx context.Context, userID string, tx *gorm.DB) ([]*Account, error) {
	db := getTxDB(am.db, tx)

	var accounts []*Account
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (am *AccountModel) GetByID(ctx context.Context, id string, tx *gorm.DB) (*Account, error) {
	db := getTxDB(am.db, tx)

	var account *Account
	if err := db.WithContext(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (am *AccountModel) GetByIDAndUserID(ctx context.Context, id string, userID string, tx *gorm.DB) (*Account, error) {
	db := getTxDB(am.db, tx)

	var account *Account
	if err := db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (am *AccountModel) GetByNameAndUserID(ctx context.Context, name string, userID string, tx *gorm.DB) (*Account, error) {
	db := getTxDB(am.db, tx)

	var account *Account
	if err := db.WithContext(ctx).Where("name = ? AND user_id = ?", name, userID).First(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}
