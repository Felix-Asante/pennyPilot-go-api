package models

import (
	"context"
	"time"

	"github.com/Felix-Asante/pennyPilot-go-api/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AllocationRule struct {
	ID        uuid.UUID                `gorm:"type:uuid;primaryKey;column:id;index"`
	TargetID  uuid.UUID                `gorm:"column:target_id;index"`
	Strategy  utils.AllocationStrategy `gorm:"column:strategy;type:allocation_strategy;not null"`
	Value     float64                  `gorm:"column:value;not null"`
	Active    bool                     `gorm:"column:active;not null;default:true"`
	CreatedAt time.Time                `gorm:"autoCreateTime"`
	UpdatedAt time.Time                `gorm:"autoUpdateTime"`
}

type AllocationRuleModel struct {
	db *gorm.DB
}

func NewAllocationRuleModel(db *gorm.DB) *AllocationRuleModel {
	return &AllocationRuleModel{db: db}
}

func (arm *AllocationRuleModel) Create(ctx context.Context, allocationRule *AllocationRule, tx *gorm.DB) error {
	db := getTxDB(arm.db, tx)

	if err := db.WithContext(ctx).Create(allocationRule).Error; err != nil {
		return err
	}

	return nil
}

func (arm *AllocationRuleModel) Save(ctx context.Context, allocationRule *AllocationRule, tx *gorm.DB) error {
	db := getTxDB(arm.db, tx)

	if err := db.WithContext(ctx).Save(allocationRule).Error; err != nil {
		return err
	}

	return nil
}

func (arm *AllocationRuleModel) Delete(ctx context.Context, allocationRule *AllocationRule, tx *gorm.DB) error {
	db := getTxDB(arm.db, tx)

	if err := db.WithContext(ctx).Delete(allocationRule).Error; err != nil {
		return err
	}

	return nil
}

func (arm *AllocationRuleModel) GetByTargetID(ctx context.Context, targetID uuid.UUID, tx *gorm.DB) ([]*AllocationRule, error) {
	db := getTxDB(arm.db, tx)

	var allocationRules []*AllocationRule
	if err := db.WithContext(ctx).Where("target_id = ?", targetID).Find(&allocationRules).Error; err != nil {
		return nil, err
	}

	return allocationRules, nil
}

func (arm *AllocationRuleModel) GetByID(ctx context.Context, id uuid.UUID, tx *gorm.DB) (*AllocationRule, error) {
	db := getTxDB(arm.db, tx)

	var allocationRule *AllocationRule
	if err := db.WithContext(ctx).Where("id = ?", id).First(&allocationRule).Error; err != nil {
		return nil, err
	}

	return allocationRule, nil
}

func (arm *AllocationRuleModel) GetByIDAndTargetID(ctx context.Context, id uuid.UUID, targetID uuid.UUID, tx *gorm.DB) (*AllocationRule, error) {
	db := getTxDB(arm.db, tx)

	var allocationRule *AllocationRule
	if err := db.WithContext(ctx).Where("id = ? AND target_id = ?", id, targetID).First(&allocationRule).Error; err != nil {
		return nil, err
	}

	return allocationRule, nil
}
